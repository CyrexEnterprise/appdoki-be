package app

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	log "github.com/sirupsen/logrus"
)

const beersTopic = "beers"
const usersTopic = "users"

type notifyService struct {
	client *messaging.Client
}

type notifier interface {
	notifyAll(topic string, content messaging.Notification)
	messageAll(topic string, content map[string]string)
}

func newNotifier(app *firebase.App) (*notifyService, error) {
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Errorf("error getting Messaging client: %v\n", err)
		return nil, err
	}

	return &notifyService{client}, nil
}

func (n *notifyService) notifyAll(topic string, content messaging.Notification) {
	n.sendMessage(&messaging.Message{
		Notification: &content,
		Topic:        topic,
	})
}

func (n *notifyService) messageAll(topic string, content map[string]string) {
	n.sendMessage(&messaging.Message{
		Data:  content,
		Topic: topic,
	})
}

func (n *notifyService) sendMessage(message *messaging.Message) {
	response, err := n.client.Send(context.Background(), message)
	if err != nil {
		log.Errorf("error sending message to topic %s: %v\n", message.Topic, err)
	}

	log.Infof("successfully sent message with id %s", response)
}
