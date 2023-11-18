package minecraft

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"pazbear-backend/cmd/minecraft/models"
	"syscall"

	amqprpc "github.com/0x4b53/amqp-rpc"
	amqprpcmw "github.com/0x4b53/amqp-rpc/middleware"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	debugLogger := log.New(os.Stdout, "DEBUG - ", log.LstdFlags)
	errorLogger := log.New(os.Stdout, "ERROR - ", log.LstdFlags)

	s := amqprpc.NewServer("amqp://admin:mk0607@localhost:5672/").
		AddMiddleware(amqprpcmw.PanicRecoveryLogging(errorLogger.Printf))

	s.WithErrorLogger(errorLogger.Printf)
	s.WithDebugLogger(debugLogger.Printf)

	s.Bind(amqprpc.TopicBinding("minecraft", "minecraft", handler))

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		s.Stop()
	}()

	s.ListenAndServe()
}

func handler(c context.Context, rw *amqprpc.ResponseWriter, d amqp091.Delivery) {
	var turnMCRequest models.Message
	err := json.Unmarshal(d.Body, &turnMCRequest)
	if err != nil {

	}
}
