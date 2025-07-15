package auth

import (
	"net/http"
	"semita/app/http/requests"
	"semita/app/http/resources"
	"semita/app/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var req requests.LoginRequest
	if err := req.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []gin.H{{
			"status": "400",
			"title":  "Validation Error",
			"detail": err.Error(),
		}}})
		return
	}

	storedUser, err := models.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": []gin.H{{
			"status": "401",
			"title":  "Unauthorized",
			"detail": "Invalid email or password",
		}}})
		return
	}

	errPassword := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(req.Password))
	if errPassword != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": []gin.H{{
			"status": "401",
			"title":  "Unauthorized",
			"detail": "Invalid email or password",
		}}})
		return
	}

	clients, err := models.GetAllClients()
	if err != nil || len(clients) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []gin.H{{
			"status": "500",
			"title":  "Server Error",
			"detail": "No OAuth client available",
		}}})
		return
	}
	client := clients[0]
	token, err := models.CreateToken(int64(storedUser.ID), client.ID, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []gin.H{{
			"status": "500",
			"title":  "Server Error",
			"detail": "Error generating OAuth token",
		}}})
		return
	}

	resource := resources.NewAuthResource(uint(storedUser.ID), storedUser.Name, storedUser.Email, token.AccessToken)
	response := resources.NewAuthLoginResponse(resource, token.RefreshToken, 86400, token.Scopes)
	c.JSON(http.StatusOK, response)
}
