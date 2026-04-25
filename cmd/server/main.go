package main

import (
	"context"
	"log"
	"os"
	"stocet-core-service/internal/config"
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/database"
	"stocet-core-service/internal/email"
	"stocet-core-service/internal/middleware"
	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/routes"
	"stocet-core-service/internal/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambdaV2

// Handler is the Lambda entry point that proxies API Gateway requests to the Gin router.
func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

// SetupRouter initialises all dependencies and returns the fully configured Gin engine.
// The function is safe to call multiple times (e.g. Lambda warm starts) because the
// underlying database connections use singleton guards.
func SetupRouter() *gin.Engine {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize PostgreSQL
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Seed initial admin user
	if err := database.SeedInitialAdmin(); err != nil {
		log.Printf("Warning: Failed to seed initial admin: %v", err)
	}

	// Initialize MongoDB
	if err := database.InitMongoDB(); err != nil {
		log.Fatal("Failed to initialize MongoDB:", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	cspRepo := repositories.NewCapitalServiceProvidersRepository()
	educationRepo := repositories.NewEducationRepository()
	chatRepo := repositories.NewChatRepository()
	alertsRepo := repositories.NewAlertsRepository()
	emailSender := email.NewSMTPSender(config.AppConfig)
	notificationsRepo := repositories.NewNotificationsRepository(emailSender)
	verificationCodeRepo := repositories.NewVerificationCodeRepository()
	newsRepo := repositories.NewNewsRepository(database.MongoDBClient)
	marketDataRepo := repositories.NewMarketDataRepository(database.DB)
	watchlistRepo := repositories.NewWatchlistRepository()

	// Initialize services
	authService := services.NewAuthService(userRepo, verificationCodeRepo, emailSender)
	cspService := services.NewCapitalServiceProvidersService(cspRepo)
	educationService := services.NewEducationService(educationRepo)
	chatService := services.NewChatService(chatRepo)
	alertsService := services.NewAlertsService(alertsRepo)
	notificationsService := services.NewNotificationsService(notificationsRepo)
	newsService := services.NewNewsService(newsRepo)
	marketDataService := services.NewMarketDataService(marketDataRepo)
	analyticsService := services.NewAnalyticsService(marketDataService)
	dashboardService := services.NewDashboardService(newsService, marketDataService, alertsService)
	watchlistService := services.NewWatchlistService(watchlistRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService, userRepo)
	cspController := controllers.NewCapitalServiceProvidersController(cspService)
	educationController := controllers.NewEducationController(educationService)
	chatController := controllers.NewChatController(chatService)
	alertsController := controllers.NewAlertsController(alertsService)
	notificationsController := controllers.NewNotificationsController(notificationsService)
	newsController := controllers.NewNewsController(newsService)
	marketDataController := controllers.NewMarketDataController(marketDataService)
	analyticsController := controllers.NewAnalyticsController(analyticsService)
	dashboardController := controllers.NewDashboardController(dashboardService)
	adminController := controllers.NewAdminController(userRepo)
	watchlistController := controllers.NewWatchlistController(watchlistService)

	// Setup Gin router
	if config.AppConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(middleware.LoggingMiddleware())
	// router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// API routes
	api := router.Group(config.API_PREFIX)
	{
		routes.SetupAuthRoutes(api, authController)
		routes.SetupNewsRoutes(api, newsController)
		routes.SetupMarketDataRoutes(api, marketDataController)
		routes.SetupCapitalServiceProviderRoutes(api, cspController)
		routes.SetupEducationRoutes(api, educationController)
		routes.SetupChatRoutes(api, chatController)
		routes.SetupAlertsRoutes(api, alertsController)
		routes.SetupNotificationsRoutes(api, notificationsController)
		routes.SetupAnalyticsRoutes(api, analyticsController)
		routes.SetupDashboardRoutes(api, dashboardController)
		routes.SetupAdminRoutes(api, adminController)
		routes.SetupWatchlistRoutes(api, watchlistController)
	}

	return router
}

func main() {
	router := SetupRouter()

	if os.Getenv("LAMBDA_TASK_ROOT") != "" {
		// Running on AWS Lambda — bridge via API Gateway adapter
		ginLambda = ginadapter.NewV2(router)
		lambda.Start(Handler)
	} else {
		// Running locally — standard HTTP server
		defer database.CloseDB()
		defer database.CloseMongoDB()

		addr := config.AppConfig.ServerHost + ":" + config.AppConfig.ServerPort
		log.Printf("Server starting on %s", addr)
		if err := router.Run(addr); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}
}
