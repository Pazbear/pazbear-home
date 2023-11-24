package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"pazbear-backend/cmd/mainapp/config"
	"pazbear-backend/cmd/mainapp/docs"

	amqprpc "github.com/0x4b53/amqp-rpc"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Controller struct {
	rpcClient *amqprpc.Client
	AppConfig config.Config
}

func NewController(appConfig config.Config) (*Controller, error) {
	rpcClient := amqprpc.NewClient(appConfig.MQ.URL)
		//ReplyToQueueName을 일정한 값으로 설정 시 apimanager가 재실행되어도 taskmanager가 보내온 반환값을 받을 수 있음
	rpcClient.WithErrorLogger(log.New(os.Stdout, "ERROR - ", log.LstdFlags).Printf)
	rpcClient.WithDebugLogger(log.New(os.Stdout, "DEBUG - ", log.LstdFlags).Printf)
	return &Controller{rpcClient:rpcClient, AppConfig: appConfig}, nil
}

// @Summary     상태 체크
// @Description	현재 서버 상태 체크
// @Tags        common
// @Router      /api/v1/healthcheck [get]
// @Success     200
func (c *Controller) HealthCheck(ginctx *gin.Context) {
	ginctx.JSON(http.StatusOK, nil)
}

func (c *Controller) NewRouter() *gin.Engine {

	r := gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/healthcheck", c.HealthCheck)

			v1.Group("/travel-log")
			{
				// v1.GET("", c.ListTravelLogs)
				// v1.GET("", c.GetTravelLog)
				// v1.POST("", c.CreateTravelLog)
				// v1.DELETE("", c.DeleteTravelLog)
			}
			mc := v1.Group("/mc")
			{
				mc.GET("", c.TurnOnServer)
				mc.GET("", c.TurnOffServer)
				mc.GET("", c.StatusServer)
			}
		}
	}

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", c.AppConfig.Address, c.AppConfig.Port)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
