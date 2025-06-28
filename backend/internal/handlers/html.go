package handlers

import (
	"brokefolio/backend/internal/middleware"
	"brokefolio/backend/internal/route"
	"log"
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
	route.RenderTemplate(w, "password-recovery.html", "")
}

// func UnauthPageHandler(w http.ResponseWriter, req *http.Request, url string, code int) {
// 	route.RenderTemplate(w, url, "")
// }

func ChangePasswordRenderHandler(w http.ResponseWriter, req *http.Request) {
	route.RenderTemplate(w, "changePassword.html", "")
}
func MarketHandler(w http.ResponseWriter, req *http.Request) {
	route.RenderTemplate(w, "market.html", "")
}
func PortfolioPageHandler(w http.ResponseWriter, req *http.Request) {
	route.RenderTemplate(w, "portfolio.html", "")
}
func MailPasswordChangeHandler(w http.ResponseWriter, req *http.Request) {

	userMail, ok := req.Context().Value(middleware.MailContext).(string)

	if !ok || userMail == "" {
		http.Error(w, "User Mail not Found", http.StatusInternalServerError)
		log.Println("User Mail not Found")
		return
	}

	route.RenderTemplate(w, "mailPasswordReset.html", userMail)
}
func UnauthPageHandler(w http.ResponseWriter, req *http.Request) {
	route.RenderTemplate(w, "unauth.html", "")
}

func WentWrongHandler(w http.ResponseWriter, req *http.Request) {
	route.RenderTemplate(w, "wentWrong.html", "")
}
