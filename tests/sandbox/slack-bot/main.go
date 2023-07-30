package main

import (
	"context"
	"log"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	client := slack.New(
		os.Getenv("SLACK_BOT_TOKEN"),
		slack.OptionAppLevelToken(os.Getenv("SLACK_APP_TOKEN")),
		slack.OptionDebug(true),
	)

	socket := socketmode.New(
		client,
		socketmode.OptionDebug(true),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context, client *slack.Client, socket *socketmode.Client) {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-socket.Events:
				switch event.Type {
				case socketmode.EventTypeEventsAPI:
					ev, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Println("ignored", event)
						continue
					}
					msg, ok := ev.InnerEvent.Data.(*slackevents.MessageEvent)
					if !ok {
						log.Println("ignored", event)
						continue
					}
					if msg.SubType == "" && msg.BotID == "" {
						_, _, err := client.PostMessage(
							msg.Channel,
							slack.MsgOptionText(msg.Text, false),
						)
						if err != nil {
							log.Printf("failed posting message: %v", err)
						}
					}
					socket.Ack(*event.Request)
				default:
					log.Println("ignored", event)
				}
			}
		}
	}(ctx, client, socket)

	if err := socket.Run(); err != nil {
		panic(err)
	}
}
