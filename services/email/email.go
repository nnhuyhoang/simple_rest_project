package email

import (
	"fmt"

	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/config"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/logger"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/util"
	"github.com/streadway/amqp"
)

// EmailService
type EmailService struct {
	Channel *amqp.Channel
	Log     logger.Log
}

func NewEmailService(cfg config.Config, l logger.Log) *EmailService {
	dialConn := fmt.Sprintf("amqp://%s:%s@%s/", cfg.RabbitUser, cfg.RabbitPassword, cfg.RabbitURL)
	conn, err := amqp.Dial(dialConn)
	util.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		"email_signup", // name
		"direct",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")
	return &EmailService{
		Channel: ch,
		Log:     l,
	}
}

func (s *EmailService) Send(exchangeName string, routeKey string, body []byte) {
	err := s.Channel.Publish(
		exchangeName, // exchange
		routeKey,     // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		s.Log.Error("[EmailService.Send] Publish()", err)
	}

}
