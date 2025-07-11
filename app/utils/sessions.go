package utils

import (
	"net/http"
	"sync"
	"web_utilidades/app/structs"

	"github.com/gorilla/sessions"
)

var sessionStoreOnce sync.Once
var sessionStore *sessions.CookieStore

func GetSessionStore() *sessions.CookieStore {
	sessionStoreOnce.Do(func() {
		appKey := GetEnv("APP_KEY")
		sessionStore = sessions.NewCookieStore([]byte(appKey))
	})
	return sessionStore
}

func LoginUserSession(response http.ResponseWriter, request *http.Request, user structs.UserStruct) error {
	var session, sessionError = GetSessionStore().Get(request, "user-session")
	if sessionError != nil {
		http.Error(response, "Error al crear la sesión", http.StatusInternalServerError)
		return sessionError
	}

	session.Values["user_id"] = user.ID
	session.Values["user_name"] = user.Name
	session.Values["user_email"] = user.Email
	session.Values["authenticated"] = true

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   84400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
	}

	return session.Save(request, response)
}

func GetAuthenticatedUser(request *http.Request) (structs.UserStruct, bool) {
	var session, sessionError = GetSessionStore().Get(request, "user-session")
	if sessionError != nil {
		return structs.UserStruct{}, false
	}

	var authenticated, ok = session.Values["authenticated"].(bool)
	if !ok || !authenticated {
		return structs.UserStruct{}, false
	}

	var user = structs.UserStruct{
		ID:    session.Values["user_id"].(int),
		Name:  session.Values["user_name"].(string),
		Email: session.Values["user_email"].(string),
	}

	return user, true
}

func LogoutUserSession(response http.ResponseWriter, request *http.Request) error {
	var session, sessionError = GetSessionStore().Get(request, "user-session")
	if sessionError != nil {
		http.Error(response, "Error al crear la sesión", http.StatusInternalServerError)
		return sessionError
	}

	session.Values["user_id"] = nil
	session.Values["user_name"] = nil
	session.Values["user_email"] = nil
	session.Values["authenticated"] = false

	session.Options.MaxAge = -1

	return session.Save(request, response)
}

func IsUserAuthenticated(request *http.Request) bool {
	_, authenticated := GetAuthenticatedUser(request)
	return authenticated
}
