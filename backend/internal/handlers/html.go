package handlers

import (
	"brokefolio/backend/internal/middleware"
	"brokefolio/backend/internal/route"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	isAuthenticated, _ := req.Context().Value(middleware.IsAuthenticatedContextKey).(bool)
	userID, _ := req.Context().Value(middleware.UserIDContextKey).(string)

	if isAuthenticated {
		route.RenderTemplate(w, "authIndex.html", userID)
	} else {
		route.RenderTemplate(w, "index.html", "")
	}
}

func LoginPageHandler(w http.ResponseWriter, req *http.Request) {
	route.RenderTemplate(w, "login.html", "")
}

func RegisterPageHandler(w http.ResponseWriter, req *http.Request) {
	route.RenderTemplate(w, "register.html", "")
}

func PassRecoveryPageHandler(w http.ResponseWriter, req *http.Request) {
	route.RenderTemplate(w, "pass-recovery.html", "")
}

// func UnauthPageHandler(w http.ResponseWriter, req *http.Request, url string, code int) {
// 	route.RenderTemplate(w, url, "")
// }

func ChangePasswordRenderHandler(w http.ResponseWriter, req *http.Request) {
	route.RenderTemplate(w, "changePassword.html", "")
}
