package controller

import (
	"encoding/json"
	"net/http"
	"pazbear-backend/cmd/mainapp/models"
	"time"

	amqprpc "github.com/0x4b53/amqp-rpc"
	"github.com/gin-gonic/gin"
)

// @Summary     서버 켜기
// @Description	마인크래프트 서버 켜기
// @Tags        mc
// @Router      /api/v1/mc/turnon [get]
// @Success     200
func (c *Controller) TurnOnServer(ginctx *gin.Context) {
	msg := models.Message{Command : "turnon"}
	msgStr, _ := json.Marshal(msg)
	res, err :=c.rpcClient.Send(amqprpc.NewRequest().WithExchange(c.AppConfig.MQ.Target.Exchange).WithRoutingKey(c.AppConfig.MQ.Target.RoutingKey).
		WithBody(string(msgStr)).WithTimeout(60 * time.Second))
	if err != nil {
		ginctx.JSON(http.StatusInternalServerError, nil)
	}
	ginctx.JSON(http.StatusOK, nil)
}

// @Summary     서버 끄기
// @Description	마인크래프트 서버 끄기
// @Tags        mc
// @Router      /api/v1/mc/turnoff [get]
// @Success     200
func (c *Controller) TurnOffServer(ginctx *gin.Context) {
	ginctx.JSON(http.StatusOK, nil)
}

// @Summary     서버 상태 확인
// @Description	마인크래프트 서버 상태 확인
// @Tags        mc
// @Router      /api/v1/mc/status [get]
// @Success     200
func (c *Controller) StatusServer(ginctx *gin.Context) {
	ginctx.JSON(http.StatusOK, nil)
}