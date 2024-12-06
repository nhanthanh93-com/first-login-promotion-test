package router

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"trinity/api/handler"
	"trinity/api/middleware"
	"trinity/api/repository"
	"trinity/api/service"
	"trinity/configs"
	"trinity/docs"
	"trinity/helpers/response"
	"trinity/pkg/app"
)

const (
	BasePath = "/api/v1"
)

func NewRouter(
	appConfig *app.Config,
	cfg *configs.Config,
) *gin.Engine {
	switch cfg.Env {
	case "prod":
		gin.SetMode(gin.ReleaseMode)
		break
	default:
		gin.SetMode(gin.DebugMode)
		break
	}
	router := gin.Default()
	cors := middleware.NewCors()
	router.Use(cors.Handler())
	router.Use(logger.SetLogger())
	router.Use(gin.Recovery())

	router.NoRoute(func(c *gin.Context) {
		response.NotFoundError(c, "Not found")
	})

	setApiGroupRoutes(router, appConfig)

	docs.SwaggerInfo.BasePath = BasePath
	ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:9888/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}

func setApiGroupRoutes(
	router *gin.Engine,
	appConfig *app.Config,
) *gin.RouterGroup {
	groups := router.Group(BasePath)

	userRepo := repository.NewUserRepository(appConfig)
	userService := service.NewUserService(appConfig, userRepo)
	userHandler := handler.NewUserHandler(userService)

	user := groups.Group("/users")
	user.POST("/create", userHandler.Create)
	user.GET("/:id", userHandler.Find)
	user.GET("/list", userHandler.List)
	user.PUT("/:id", userHandler.Update)
	user.DELETE("/:id", userHandler.Delete)

	promo := groups.Group("/promo")
	promo.POST("/register", userHandler.Register)

	campaignRepo := repository.NewCampaignRepository(appConfig)
	campaignService := service.NewCampaignService(campaignRepo)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	campaign := groups.Group("/campaigns")

	campaign.POST("/create", campaignHandler.Create)
	campaign.GET("/:id", campaignHandler.Find)
	campaign.GET("/list", campaignHandler.List)
	campaign.PUT("/:id", campaignHandler.Update)
	campaign.DELETE("/:id", campaignHandler.Delete)

	voucherRepo := repository.NewVoucherRepository(appConfig)
	voucherService := service.NewVoucherService(voucherRepo)
	voucherHandler := handler.NewVoucherHandler(voucherService)

	voucher := groups.Group("/vouchers")
	voucher.GET("/:id", voucherHandler.Find)
	voucher.GET("/list", voucherHandler.List)
	voucher.PUT("/:id", voucherHandler.Update)
	voucher.DELETE("/:id", voucherHandler.Delete)

	productRepo := repository.NewProductRepository(appConfig)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	product := groups.Group("/products")
	product.POST("/create", productHandler.Create)
	product.GET("/:id", productHandler.Find)
	product.GET("/list", productHandler.List)
	product.PUT("/:id", productHandler.Update)
	product.DELETE("/:id", productHandler.Delete)

	cartRepo := repository.NewCartRepository(appConfig)
	cartService := service.NewCartService(appConfig, cartRepo)
	cartHandler := handler.NewCartHandler(cartService)

	cart := groups.Group("/carts")
	cart.GET("/:id", cartHandler.Find)
	cart.DELETE("/id", cartHandler.DeleteCartItem)

	order := groups.Group("/order")
	order.POST("/:user_id", cartHandler.CreateOrder)
	order.PUT("/:id/status", cartHandler.UpdateOrderStatus)

	return groups
}
