package main

import (
	"brokefolio/backend/internal/database"
	"brokefolio/backend/internal/handlers"
	"brokefolio/backend/internal/middleware"
	"brokefolio/backend/internal/route"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/cors"
)

func main() {

	route.InitTemplates()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	staticPath := filepath.Join(wd, "../../../static")

	db := database.InitDB()
	defer db.Close()

	router := http.NewServeMux()

	registerHandler := handlers.NewRegisterHandler(db)
	loginHandler := handlers.NewLoginHandler(db)
	mailHandler := handlers.NewMailHandler(db)
	sessionHandler := middleware.NewSessionMiddleware(db)

	profileHandler := handlers.NewProfileHandler(db)

	fs := http.FileServer(http.Dir(staticPath))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	router.HandleFunc("POST /api/login", loginHandler.LoginHandler)
	router.HandleFunc("POST /api/register", registerHandler.RegisterHandler)
	router.HandleFunc("POST /api/forgetPassword", mailHandler.MailResetHandler)
	router.Handle("/homepage", sessionHandler.CheckSessionMiddleware(http.HandlerFunc(handlers.HomeHandler)))
	//market
	//portfolio
	router.Handle("/profile", sessionHandler.CheckSessionMiddleware(http.HandlerFunc(profileHandler.ProfilePageHandler)))
	router.HandleFunc("/login", handlers.LoginPageHandler)
	router.HandleFunc("/register", handlers.RegisterPageHandler)
	router.HandleFunc("/passrecover", handlers.PassRecoveryPageHandler)
	//ChangePasswordRenderHandler
	router.Handle("/change-password", sessionHandler.CheckSessionMiddleware(http.HandlerFunc(handlers.ChangePasswordRenderHandler)))
	//user-apis
	router.Handle("PUT /api/user/update-profile", middleware.MiddlewareAuthJWT(http.Handler(sessionHandler.CheckSessionMiddleware(http.HandlerFunc(profileHandler.UpdateProfileHandler)))))
	router.Handle("DELETE /api/user/delete-account", middleware.MiddlewareAuthJWT(http.Handler(sessionHandler.CheckSessionMiddleware(http.HandlerFunc(profileHandler.DeleteProfileHandler)))))
	router.Handle("POST /api/user/change-password", middleware.MiddlewareAuthJWT(http.Handler(sessionHandler.CheckSessionMiddleware(http.HandlerFunc(profileHandler.ChangePasswordHandler)))))

	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsStr == "" {
		allowedOriginsStr = "http://127.0.0.1:8080,http://localhost:8080"
	}

	allowedOrigins := strings.Split(allowedOriginsStr, ",")

	c := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		Debug:          os.Getenv("CORS_DEBUG") == "true",
	})

	handler := c.Handler(router)

	log.Println("Server listening on :8080")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
