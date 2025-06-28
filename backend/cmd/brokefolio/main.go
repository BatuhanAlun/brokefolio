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

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	if err := godotenv.Load("../../../.env"); err != nil {
		log.Println("No .env file found, assuming environment variables are set externally.")
	}

	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")

	route.InitTemplates()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	staticPath := filepath.Join(wd, "../../../static")

	db := database.InitDB(dbConnectionString)
	defer db.Close()

	router := http.NewServeMux()

	registerHandler := handlers.NewRegisterHandler(db)
	loginHandler := handlers.NewLoginHandler(db)
	mailHandler := handlers.NewMailHandler(db)
	sessionHandler := middleware.NewSessionMiddleware(db)
	profileHandler := handlers.NewProfileHandler(db)
	tradeHandler := handlers.NewBuyHandler(db)
	priceFetcher := &handlers.HTTPPriceFetcher{}
	portfolioHandler := handlers.NewPortfolioHandler(db, priceFetcher)
	transactionHandler := handlers.NewTransactionsHandler(db)
	mailMiddlewareHandler := middleware.NewMailMiddlewareDB(db)
	logoutHandler := handlers.NewLogOutHandler(db)

	// Static files
	fs := http.FileServer(http.Dir(staticPath))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	// Public API routes
	router.HandleFunc("POST /api/login", loginHandler.LoginHandler)
	router.HandleFunc("POST /api/register", registerHandler.RegisterHandler)
	router.HandleFunc("POST /api/forgetPassword", mailHandler.MailResetHandler)
	router.HandleFunc("POST /api/user/change-password-mail", mailHandler.PasswordResetUsingMailHandler)
	router.HandleFunc("GET /api/news", handlers.CombinedNewsHandler)
	router.HandleFunc("GET /api/crypto-price", handlers.StockPriceHandler)

	// Public page routes
	router.HandleFunc("/login", handlers.LoginPageHandler)
	router.HandleFunc("/register", handlers.RegisterPageHandler)
	router.HandleFunc("/passrecover", handlers.PassRecoveryPageHandler)

	// Password reset via email link
	router.Handle("/reset-password", mailMiddlewareHandler.CheckMailParameterMiddleware(http.HandlerFunc(handlers.MailPasswordChangeHandler)))

	// Protected trade route
	router.Handle("POST /api/trade", sessionHandler.CheckSessionMiddleware(http.HandlerFunc(tradeHandler.TradeHandler)))

	// Authenticated API routes
	router.Handle("PUT /api/user/update-profile",
		middleware.MiddlewareAuthJWT(http.Handler(
			sessionHandler.CheckSessionMiddleware(http.HandlerFunc(profileHandler.UpdateProfileHandler))),
		),
	)
	router.Handle("DELETE /api/user/delete-account",
		middleware.MiddlewareAuthJWT(http.Handler(
			sessionHandler.CheckSessionMiddleware(http.HandlerFunc(profileHandler.DeleteProfileHandler))),
		),
	)
	router.Handle("POST /api/user/change-password",
		middleware.MiddlewareAuthJWT(http.Handler(
			sessionHandler.CheckSessionMiddleware(http.HandlerFunc(profileHandler.ChangePasswordHandler))),
		),
	)
	router.Handle("/api/portfolio",
		middleware.MiddlewareAuthJWT(http.Handler(
			sessionHandler.CheckSessionMiddleware(http.HandlerFunc(portfolioHandler.PortfolioHandler))),
		),
	)
	router.Handle("/api/transactions",
		middleware.MiddlewareAuthJWT(http.Handler(
			sessionHandler.CheckSessionMiddleware(http.HandlerFunc(transactionHandler.TransactionsHandler))),
		),
	)

	// Session/middleware-protected pages
	router.Handle("/homepage", sessionHandler.CheckSessionMiddleware(http.HandlerFunc(handlers.HomeHandler)))
	router.Handle("/market", sessionHandler.CheckSessionMiddleware(http.HandlerFunc(handlers.MarketHandler)))
	router.Handle("/portfolio", sessionHandler.CheckSessionMiddleware(http.HandlerFunc(handlers.PortfolioPageHandler)))
	router.Handle("/profile", sessionHandler.CheckSessionMiddleware(http.HandlerFunc(profileHandler.ProfilePageHandler)))
	router.Handle("/change-password", sessionHandler.CheckSessionMiddleware(http.HandlerFunc(handlers.ChangePasswordRenderHandler)))

	// Logout
	router.HandleFunc("/logout", logoutHandler.LogoutHandler)

	// Error/Fallback pages
	router.HandleFunc("/unauthorized", handlers.UnauthPageHandler)
	router.HandleFunc("/went-wrong", handlers.WentWrongHandler)

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
