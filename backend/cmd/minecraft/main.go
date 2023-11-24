package minecraft

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"pazbear-backend/cmd/minecraft/models"
	"syscall"

	amqprpc "github.com/0x4b53/amqp-rpc/v3"
	amqprpcmw "github.com/0x4b53/amqp-rpc/v3/middleware"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	debugLogger := log.New(os.Stdout, "DEBUG - ", log.LstdFlags)
	errorLogger := log.New(os.Stdout, "ERROR - ", log.LstdFlags)

	s := amqprpc.NewServer("amqp://admin:mk0607@localhost:5672/").
		AddMiddleware(amqprpcmw.PanicRecoveryLogging(errorLogger.Printf))

	s.WithErrorLogger(errorLogger.Printf)
	s.WithDebugLogger(debugLogger.Printf)

	s.Bind(amqprpc.TopicBinding("minecraft", "minecraft", handle))

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
	switch msg.Command {
	case "turnon":
		turnonCmd := exec.Command("docker-compose", "-f", "/Users/mkcho1997/Project/pazbear-home/minecraft/dawncraft-docker-compose.yml", "up", "-d")
		turnonCmd.Stdout = os.Stdout
		turnonCmd.Stderr = os.Stderr
		if err := turnonCmd.Run(); err != nil {
			res.Success = false
			res.Output = models.ErrorLog{
				Err: err.Error(),
			}
			fmt.Fprint(rw, res)
		}else{
			res.Success = true
			res.Output = models.OutputLog{
				Output: "Minecraft Server is Started",
			}
			fmt.Fprint(rw, res)
		}

	case "turnoff":
		turnoffCmd := exec.Command("docker-compose", "exec", "-i", "minecraft-mc-1", "rcon-cli", "stop")
		turnoffCmd.Stdout = os.Stdout
		turnoffCmd.Stderr = os.Stderr
		if err := turnoffCmd.Run(); err != nil {
			res.Success = false
			res.Output = models.ErrorLog{
				Err: err.Error(),
			}
			fmt.Fprint(rw, res)
		}else{
			res.Success = true
			res.Output = models.OutputLog{
				Output: "Minecraft Server is Stopped",
			}
			fmt.Fprint(rw, res)
		}
	case "status":
		
	}
	d.Ack(false)
}
