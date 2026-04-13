package main

import (
	"log"
	"stocet-core-service/internal/config"
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/database"
	"stocet-core-service/internal/middleware"
	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/routes"
	"stocet-core-service/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize PostgreSQL
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Seed initial admin user
	if err := database.SeedInitialAdmin(); err != nil {
		log.Printf("Warning: Failed to seed initial admin: %v", err)
	}

	// Initialize MongoDB
	if err := database.InitMongoDB(); err != nil {
		log.Fatal("Failed to initialize MongoDB:", err)
	}
	defer database.CloseMongoDB()

	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	brokersRepo := repositories.NewBrokersRepository()
	educationRepo := repositories.NewEducationRepository()
	chatRepo := repositories.NewChatRepository()
	alertsRepo := repositories.NewAlertsRepository()
	newsRepo := repositories.NewNewsRepository(database.MongoDBClient)
	marketDataRepo := repositories.NewMarketDataRepository(database.MongoDBClient)
	watchlistRepo := repositories.NewWatchlistRepository()

	// Initialize services
	authService := services.NewAuthService(userRepo)
	brokersService := services.NewBrokersService(brokersRepo)
	educationService := services.NewEducationService(educationRepo)
	chatService := services.NewChatService(chatRepo)
	alertsService := services.NewAlertsService(alertsRepo)
	newsService := services.NewNewsService(newsRepo)
	marketDataService := services.NewMarketDataService(marketDataRepo)
	analyticsService := services.NewAnalyticsService(marketDataService)
	dashboardService := services.NewDashboardService(newsService, marketDataService, alertsService)
	watchlistService := services.NewWatchlistService(watchlistRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService, userRepo)
	brokersController := controllers.NewBrokersController(brokersService)
	educationController := controllers.NewEducationController(educationService)
	chatController := controllers.NewChatController(chatService)
	alertsController := controllers.NewAlertsController(alertsService)
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
	router.Use(middleware.CORSMiddleware())
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
		routes.SetupBrokersRoutes(api, brokersController)
		routes.SetupEducationRoutes(api, educationController)
		routes.SetupChatRoutes(api, chatController)
		routes.SetupAlertsRoutes(api, alertsController)
		routes.SetupAnalyticsRoutes(api, analyticsController)
		routes.SetupDashboardRoutes(api, dashboardController)
		routes.SetupAdminRoutes(api, adminController)
		routes.SetupWatchlistRoutes(api, watchlistController)
	}

	// Start server
	addr := config.AppConfig.ServerHost + ":" + config.AppConfig.ServerPort
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

