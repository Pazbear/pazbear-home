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
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		fmt.Fprint(rw, "invalid message type")
	}
	switch msg.Command {
	case "turnon":
		turnonCmd := exec.Command("docker-compose", "-f", "/Users/mkcho1997/Project/pazbear-home/minecraft/dawncraft-docker-compose.yml", "up", "-d")
		turnonCmd.Stdout = os.Stdout
		turnonCmd.Stderr = os.Stderr
		if err := turnonCmd.Run(); err != nil {
			fmt.Println(err)
		}
		d.Ack(false)

	case "turnoff":
		turnonCmd := exec.Command("docker-compose", "exec", "-i", "minecraft-mc-1", "rcon-cli", "stop")
		turnonCmd.Stdout = os.Stdout
		turnonCmd.Stderr = os.Stderr
		if err := turnonCmd.Run(); err != nil {
			fmt.Println(err)
		}
		d.Ack(false)
	case "status":

	}
}
