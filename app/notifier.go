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
	client           *messaging.Client
	androidMsgConfig *messaging.AndroidConfig
	apnsMsgConfig    *messaging.APNSConfig
}

type notifier interface {
	notifyAll(topic string, notification *messaging.Notification, data map[string]string)
	messageAll(topic string, content map[string]string)
}

func newNotifier(app *firebase.App) (*notifyService, error) {
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Errorf("error getting Messaging client: %v\n", err)
		return nil, err
	}

	return &notifyService{
		client: client,
		androidMsgConfig: &messaging.AndroidConfig{
			Priority: "high",
		},
		apnsMsgConfig: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					ContentAvailable: true,
				},
			},
		},
	}, nil
}

func (n *notifyService) notifyAll(topic string, notification *messaging.Notification, data map[string]string) {
	n.sendMessage(&messaging.Message{
		Data:         data,
		Notification: notification,
		Topic:        topic,
		Android:      n.androidMsgConfig,
		APNS:         n.apnsMsgConfig,
	})
}

func (n *notifyService) messageAll(topic string, content map[string]string) {
	n.sendMessage(&messaging.Message{
		Data:    content,
		Topic:   topic,
		Android: n.androidMsgConfig,
		APNS:    n.apnsMsgConfig,
	})
}

func (n *notifyService) sendMessage(message *messaging.Message) {
	response, err := n.client.Send(context.Background(), message)
	if err != nil {
		log.Errorf("error sending message to topic %s: %v\n", message.Topic, err)
	}

	log.Infof("successfully sent message with id %s", response)
}
