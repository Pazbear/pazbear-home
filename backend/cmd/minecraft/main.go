package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"pazbear-backend/cmd/minecraft/config"
	"pazbear-backend/cmd/minecraft/models"
	"syscall"

	amqprpc "github.com/0x4b53/amqp-rpc/v3"
	amqprpcmw "github.com/0x4b53/amqp-rpc/v3/middleware"
	amqp "github.com/rabbitmq/amqp091-go"
)

var appConfig config.Config

func main() {
	appConfig, err := config.AppConfig()
	if err != nil {
		log.Panic(err)
	}

	debugLogger := log.New(os.Stdout, "DEBUG - ", log.LstdFlags)
	errorLogger := log.New(os.Stdout, "ERROR - ", log.LstdFlags)

	s := amqprpc.NewServer(appConfig.MQ.URL).
		AddMiddleware(amqprpcmw.PanicRecoveryLogging(errorLogger.Printf))

	s.WithErrorLogger(errorLogger.Printf)
	s.WithDebugLogger(debugLogger.Printf)

	binding := amqprpc.TopicBinding(appConfig.MQ.Target.RoutingKey, appConfig.MQ.Target.RoutingKey, handle)
	binding.ExchangeName = appConfig.MQ.Target.Exchange
	s.Bind(binding)

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		s.Stop()
	}()

	s.ListenAndServe()
}

func handle(c context.Context, rw *amqprpc.ResponseWriter, d amqp.Delivery) {
	var msg models.Message
	var res models.MQResponse
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		res.Success = false
		res.Output = models.ErrorLog{
			Err: "invalid message type",
		}
		fmt.Fprint(rw, res)
	}
	fmt.Println(msg.Command)
	switch msg.Command {
	case "turnon":
		turnonCmd := exec.Command("docker-compose", "-f", "/Users/mkcho1997/Project/pazbear-home/minecraft/darkrpg-docker-compose.yml", "up", "-d")
		turnonCmd.Stdout = os.Stdout
		turnonCmd.Stderr = os.Stderr
		if err := turnonCmd.Run(); err != nil {
			res.Success = false
			res.Output = models.ErrorLog{
				Err: err.Error(),
			}
			result, _ := json.Marshal(res)
			fmt.Fprint(rw, string(result))
		} else {
			res.Success = true
			res.Output = models.OutputLog{
				Output: "Minecraft Server is Started",
			}
			result, _ := json.Marshal(res)
			fmt.Fprint(rw, string(result))
		}

	case "turnoff":
		turnoffCmd := exec.Command("docker", "exec", "-i", "minecraft-mc-1", "rcon-cli", "stop")
		turnoffCmd.Stdout = os.Stdout
		turnoffCmd.Stderr = os.Stderr
		if err := turnoffCmd.Run(); err != nil {
			res.Success = false
			res.Output = models.ErrorLog{
				Err: err.Error(),
			}
			result, _ := json.Marshal(res)
			fmt.Fprint(rw, string(result))
		} else {
			res.Success = true
			res.Output = models.OutputLog{
				Output: "Minecraft Server is Stopped",
			}
			result, _ := json.Marshal(res)
			fmt.Fprint(rw, string(result))
		}
	case "status":
		statusCmd := exec.Command("docker", "inspect", "-f", "'{{.State.Status}}'", "minecraft-mc-1")
		if statusCmdOut, err := statusCmd.Output(); err != nil {
			res.Success = false
			res.Output = models.ErrorLog{
				Err: err.Error(),
			}
			result, _ := json.Marshal(res)
			fmt.Fprint(rw, string(result))
		} else {
			res.Success = true
			res.Output = models.SvrStatus{
				Status: string(statusCmdOut),
			}
			result, _ := json.Marshal(res)
			fmt.Fprint(rw, string(result))
		}
	}
}
