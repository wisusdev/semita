package helpers

import (
	"net/http"
	"web_utilidades/app/structs"
	"web_utilidades/app/utils"
)

func AuthSessionService(response http.ResponseWriter, request *http.Request, title string, data interface{}) structs.AuthSessionStruct {
	user, isAuthenticated := utils.GetAuthenticatedUser(request)
	alertId, alertMessage := utils.GetFlashNotifications(response, request)

	lang := "es"
	if cookie, err := request.Cookie("lang"); err == nil {
		lang = cookie.Value
	}

	translate := func(key string) string {
		return utils.Translate(key, lang)
	}

	return structs.AuthSessionStruct{
		User:            user,
		IsAuthenticated: isAuthenticated,
		Title:           title,
		Data:            data,
		AlertId:         alertId,
		AlertMessage:    alertMessage,
		Lang:            lang,
		Translate:       translate,
	}
}
