package api

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"trinity/api/handler"
	"trinity/api/middleware"
	"trinity/api/repository"
	"trinity/api/service"
	"trinity/cmd/docs"
	"trinity/configs"
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

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pongs",
	})
}

func setApiGroupRoutes(
	router *gin.Engine,
	appConfig *app.Config,
) *gin.RouterGroup {
	groups := router.Group(BasePath)

	pings := groups.Group("/ping")
	pings.GET("", Ping)

	userRepo := repository.NewUserRepository(appConfig)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	user := groups.Group("/users")
	user.POST("/register", userHandler.Register)
	user.POST("/create", userHandler.Create)
	user.GET("/:id", userHandler.Find)
	user.GET("/list", userHandler.List)
	user.PUT("/:id", userHandler.Update)
	user.DELETE("/:id", userHandler.Delete)

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

	return groups
}
