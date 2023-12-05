package controller

import (
	"encoding/json"
	"net/http"
	"pazbear-backend/cmd/mainapp/models"
	"time"

	amqprpc "github.com/0x4b53/amqp-rpc"
	"github.com/gin-gonic/gin"
)

// @Summary     서버 켜기/끄기
// @Description	마인크래프트 서버 켜기/끄기
// @Tags        mc
// @Accept      json
// @Produce     json
// @Param       request       body   models.TurnMCRequest true "서버 전원 ON/OFF"
// @Router      /api/v1/mc/turn [post]
// @Success     200
func (c *Controller) TurnServer(ginctx *gin.Context) {
	var req models.TurnMCRequest
	if err := ginctx.Bind(&req); err != nil {
		ginctx.JSON(http.StatusBadRequest, models.JSONError{Error: "invalid request"})
		return
	}
	var msg models.Message
	if req.TurnOnOff {
		msg = models.Message{Command: "turnon"}
	} else {
		msg = models.Message{Command: "turnoff"}
	}

	msgStr, _ := json.Marshal(msg)
	res, err := c.rpcClient.Send(amqprpc.NewRequest().WithExchange(c.AppConfig.MQ.Target.Exchange).WithRoutingKey(c.AppConfig.MQ.Target.RoutingKey).
		WithBody(string(msgStr)).WithTimeout(60 * time.Second))
	if err != nil {
		ginctx.JSON(http.StatusInternalServerError, models.JSONError{Error: err.Error()})
		return
	}
	var mqRes models.MQResponse
	mqRes.Output = ""
	err = json.Unmarshal(res.Body, &mqRes)
	if err != nil {
		ginctx.JSON(http.StatusInternalServerError, models.JSONError{Error: err.Error()})
		return
	}
	ginctx.JSON(http.StatusOK, mqRes)
}

// @Summary     서버 상태 확인
// @Description	마인크래프트 서버 상태 확인
// @Tags        mc
// @Router      /api/v1/mc/status [get]
// @Success     200
func (c *Controller) StatusServer(ginctx *gin.Context) {
	msg := models.Message{Command: "status"}
	msgStr, _ := json.Marshal(msg)
	res, err := c.rpcClient.Send(amqprpc.NewRequest().WithExchange(c.AppConfig.MQ.Target.Exchange).WithRoutingKey(c.AppConfig.MQ.Target.RoutingKey).
		WithBody(string(msgStr)).WithTimeout(60 * time.Second))
	if err != nil {
		ginctx.JSON(http.StatusInternalServerError, models.JSONError{Error: err.Error()})
		return
	}
	var mqRes models.MQResponse
	mqRes.Output = models.SvrStatus{}
	err = json.Unmarshal(res.Body, &mqRes)
	if err != nil {
		ginctx.JSON(http.StatusInternalServerError, models.JSONError{Error: err.Error()})
		return
	}
	ginctx.JSON(http.StatusOK, mqRes)
}
