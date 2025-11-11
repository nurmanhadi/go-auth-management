package publisher

import (
	"auth-management/internal/event"
	"context"
	"time"

	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type UserPublisher struct {
	ch *amqp.Channel
}

func NewUserPublisher(ch *amqp.Channel) *UserPublisher {
	return &UserPublisher{
		ch: ch,
	}
}
func (e *UserPublisher) PublishUserRegistered(data *event.UserRegisteredPublish) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = e.ch.PublishWithContext(ctx,
		"auth.management",
		"user.registered",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
