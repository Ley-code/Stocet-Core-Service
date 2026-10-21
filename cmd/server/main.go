package main

import (
	"context"
	"log"
	"time"

	"stocet-core-service/internal/config"
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/database"
	"stocet-core-service/internal/email"
	"stocet-core-service/internal/middleware"
	"stocet-core-service/internal/push"
	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/routes"
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/services/ai"

	"github.com/gin-gonic/gin"
)

// SetupRouter initialises all dependencies and returns the fully configured Gin engine.
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

	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	middleware.SetAdminUserRepository(userRepo)
	cspRepo := repositories.NewCapitalServiceProvidersRepository()
	educationRepo := repositories.NewEducationRepository()
	chatRepo := repositories.NewChatRepository()
	alertsRepo := repositories.NewAlertsRepository()
	emailSender := email.NewSMTPSender(config.AppConfig)
	fcmService := push.NewFCMService(config.AppConfig.FCMProjectID, config.AppConfig.FCMServiceAccountJSON)
	notificationsRepo := repositories.NewNotificationsRepository(emailSender, fcmService)
	verificationCodeRepo := repositories.NewVerificationCodeRepository()
	inviteCodeRepo := repositories.NewInviteCodeRepository()
	orgRepo := repositories.NewOrganizationRepository()
	newsRepo := repositories.NewPostgresNewsRepository(database.DB)
	articleBookmarksRepo := repositories.NewPostgresArticleBookmarksRepository(database.DB)
	articleInteractionsRepo := repositories.NewPostgresArticleInteractionsRepository(database.DB)
	marketDataRepo := repositories.NewMarketDataRepository(database.DB)
	chartLayoutRepo := repositories.NewChartLayoutRepository(database.DB)
	paperTradeRepo := repositories.NewPaperTradeRepository(database.DB)
	watchlistRepo := repositories.NewWatchlistRepository()
	portfolioRepo := repositories.NewPortfolioRepository()
	aiQuotaRepo := repositories.NewAIQuotaRepository()
	fundamentalsRepo := repositories.NewFundamentalAnalysisRepository()

	// Initialize AI / Gemini stack first so other services can opt in to embedding hooks.
	geminiClient := ai.NewGeminiClient(
		config.AppConfig.GeminiAPIKey,
		config.AppConfig.GeminiChatModel,
		config.AppConfig.GeminiEmbedModel,
		time.Duration(config.AppConfig.AIRequestTimeoutSec)*time.Second,
	)
	embedder := ai.NewEmbedder(database.DB, geminiClient)
	ragService := ai.NewRAGService(database.DB, embedder, fundamentalsRepo)
	taService := ai.NewTechnicalAnalysisService(marketDataRepo)
	aiService := services.NewAIService(
		geminiClient,
		ragService,
		taService,
		config.AppConfig.AIChatHistoryMaxTurns,
		config.AppConfig.AIChatHistoryMaxChars,
	)

	// Initialize services
	subscriptionService := services.NewSubscriptionService(userRepo, orgRepo)
	authService := services.NewAuthService(userRepo, verificationCodeRepo, inviteCodeRepo, orgRepo, subscriptionService, emailSender)
	orgService := services.NewOrganizationService(orgRepo, userRepo)
	cspService := services.NewCapitalServiceProvidersService(cspRepo)
	educationService := services.NewEducationService(educationRepo)
	chatService := services.NewChatService(chatRepo)
	alertsService := services.NewAlertsService(alertsRepo, subscriptionService)
	notificationsService := services.NewNotificationsService(notificationsRepo)
	newsService := services.NewNewsService(newsRepo)
	articleBookmarksService := services.NewArticleBookmarksService(articleBookmarksRepo)
	articleInteractionsService := services.NewArticleInteractionsService(articleInteractionsRepo)
	marketDataService := services.NewMarketDataService(marketDataRepo, taService)
	chartLayoutService := services.NewChartLayoutService(chartLayoutRepo)
	paperTradeController := controllers.NewPaperTradeController(services.NewPaperTradeService(paperTradeRepo))
	analyticsService := services.NewAnalyticsService(marketDataService)
	dashboardService := services.NewDashboardService(newsService, marketDataService, alertsService)
	watchlistService := services.NewWatchlistService(watchlistRepo)
	portfolioService := services.NewPortfolioService(portfolioRepo, marketDataService)
	aiQuotaService := services.NewAIQuotaService(aiQuotaRepo, subscriptionService)
	fundamentalsService := services.NewFundamentalsService(fundamentalsRepo).WithEmbedder(embedder)

	// Initialize controllers
	authController := controllers.NewAuthController(authService, userRepo, subscriptionService)
	cspController := controllers.NewCapitalServiceProvidersController(cspService)
	educationController := controllers.NewEducationController(educationService)
	chatController := controllers.NewChatController(chatService)
	alertsController := controllers.NewAlertsController(alertsService)
	notificationsController := controllers.NewNotificationsController(notificationsService)
	notifPrefRepo := repositories.NewNotificationPreferenceRepository()
	notifPrefController := controllers.NewNotificationPreferenceController(notifPrefRepo)
	newsController := controllers.NewNewsController(
		newsService,
		articleInteractionsService,
		articleBookmarksService,
		watchlistService,
		portfolioService,
	)
	articleInteractionsController := controllers.NewArticleInteractionsController(articleInteractionsService)
	marketDataController := controllers.NewMarketDataController(marketDataService, taService, chartLayoutService)
	analyticsController := controllers.NewAnalyticsController(analyticsService)
	dashboardController := controllers.NewDashboardController(dashboardService)
	adminController := controllers.NewAdminController(userRepo)
	inviteCodesAdminController := controllers.NewInviteCodesAdminController(inviteCodeRepo)
	organizationController := controllers.NewOrganizationController(orgService)
	fundamentalsController := controllers.NewFundamentalsController(fundamentalsService, services.NewFinancialImportService(fundamentalsRepo))
	embeddingsAdminController := controllers.NewEmbeddingsAdminController(embedder).WithSentiment(geminiClient, newsRepo)
	watchlistController := controllers.NewWatchlistController(watchlistService)
	portfolioController := controllers.NewPortfolioController(portfolioService)
	advisoryController := controllers.NewAdvisoryController(repositories.NewAdvisoryReportRepository())
	screenerController := controllers.NewScreenerController()
	caRepo := repositories.NewCorporateActionRepository()
	corporateActionController := controllers.NewCorporateActionController(caRepo)
	earningsRepo := repositories.NewEarningsRepository()
	earningsController := controllers.NewEarningsController(earningsRepo)
	auditRepo := repositories.NewAuditLogRepository()
	auditController := controllers.NewAuditLogController(auditRepo)
	adminAnalyticsController := controllers.NewAdminAnalyticsController()
	aiController := controllers.NewAIController(aiService, aiQuotaService, userRepo)

	// Setup Gin router
	if config.AppConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(middleware.LoggingMiddleware())
	// The frontend is on another origin, so the API itself must answer CORS headers and preflight requests.
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})
	betaController := controllers.NewBetaController(emailSender)
	router.GET("/health/data", betaController.DataHealth)

	// API routes
	api := router.Group(config.API_PREFIX)
	{
		routes.SetupAuthRoutes(api, authController)
		routes.SetupNewsRoutes(api, newsController, articleInteractionsController)
		routes.SetupMarketDataRoutes(api, marketDataController)
		routes.SetupPaperTradeRoutes(api, paperTradeController)
		routes.SetupBetaRoutes(api, betaController)
		routes.SetupCapitalServiceProviderRoutes(api, cspController)
		routes.SetupEducationRoutes(api, educationController)
		routes.SetupChatRoutes(api, chatController)
		routes.SetupAlertsRoutes(api, alertsController)
		routes.SetupNotificationsRoutes(api, notificationsController, notifPrefController)
		routes.SetupAnalyticsRoutes(api, analyticsController)
		routes.SetupDashboardRoutes(api, dashboardController)
		marketDataAdminController := controllers.NewMarketDataAdminController()
		routes.SetupAdminRoutes(api, adminController, inviteCodesAdminController, organizationController, corporateActionController, auditController, adminAnalyticsController, fundamentalsController, embeddingsAdminController, marketDataAdminController, chatController, auditRepo, earningsController)
		routes.SetupFundamentalsPublicRoutes(api, fundamentalsController)
		routes.SetupMarketOverviewRoutes(api, controllers.NewMarketOverviewController())
		routes.SetupWatchlistRoutes(api, watchlistController)
		routes.SetupPortfolioRoutes(api, portfolioController, subscriptionService)
		routes.SetupAdvisoryRoutes(api, advisoryController, subscriptionService)
		routes.SetupScreenerRoutes(api, screenerController)
		routes.SetupCorporateActionRoutes(api, corporateActionController)
		routes.SetupEarningsRoutes(api, earningsController)
		routes.SetupAIRoutes(api, aiController, subscriptionService)
		routes.SetupOrgRoutes(api, organizationController)
	}

	// Background workers: fill in missing embeddings and classify the sentiment of new articles.
	go ai.StartBackfillWorker(
		context.Background(),
		embedder,
		time.Duration(config.AppConfig.AIBackfillIntervalSec)*time.Second,
		config.AppConfig.AIBackfillBatchSize,
	)
	go ai.StartSentimentWorker(
		context.Background(),
		geminiClient,
		newsRepo,
		15*time.Minute,
		20,
	)

	return router
}

func main() {
	router := SetupRouter()
	defer database.CloseDB()

	addr := config.AppConfig.ServerHost + ":" + config.AppConfig.ServerPort
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
