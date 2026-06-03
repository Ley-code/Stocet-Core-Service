package main

import (
	"context"
	"log"
	"os"
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
	articleInteractionsRepo := repositories.NewPostgresArticleInteractionsRepository(database.DB)
	marketDataRepo := repositories.NewMarketDataRepository(database.DB)
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
	articleInteractionsService := services.NewArticleInteractionsService(articleInteractionsRepo)
	marketDataService := services.NewMarketDataService(marketDataRepo)
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
	newsController := controllers.NewNewsController(newsService, articleInteractionsService)
	articleInteractionsController := controllers.NewArticleInteractionsController(articleInteractionsService)
	marketDataController := controllers.NewMarketDataController(marketDataService)
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
		routes.SetupNewsRoutes(api, newsController, articleInteractionsController)
		routes.SetupMarketDataRoutes(api, marketDataController)
		routes.SetupCapitalServiceProviderRoutes(api, cspController)
		routes.SetupEducationRoutes(api, educationController)
		routes.SetupChatRoutes(api, chatController)
		routes.SetupAlertsRoutes(api, alertsController)
		routes.SetupNotificationsRoutes(api, notificationsController, notifPrefController)
		routes.SetupAnalyticsRoutes(api, analyticsController)
		routes.SetupDashboardRoutes(api, dashboardController)
		routes.SetupAdminRoutes(api, adminController, inviteCodesAdminController, organizationController, corporateActionController, auditController, adminAnalyticsController, fundamentalsController, embeddingsAdminController, auditRepo)
		routes.SetupFundamentalsPublicRoutes(api, fundamentalsController)
		routes.SetupWatchlistRoutes(api, watchlistController)
		routes.SetupPortfolioRoutes(api, portfolioController, subscriptionService)
		routes.SetupAdvisoryRoutes(api, advisoryController, subscriptionService)
		routes.SetupScreenerRoutes(api, screenerController)
		routes.SetupCorporateActionRoutes(api, corporateActionController)
		routes.SetupAIRoutes(api, aiController, subscriptionService)
	}

	// Long-running embedding backfill is only safe outside Lambda (where
	// the process is recycled per invocation). On Lambda, schedule a
	// separate job to call ai.RunBackfillOnce(...) periodically.
	if os.Getenv("LAMBDA_TASK_ROOT") == "" {
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

		addr := config.AppConfig.ServerHost + ":" + config.AppConfig.ServerPort
		log.Printf("Server starting on %s", addr)
		if err := router.Run(addr); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}
}
