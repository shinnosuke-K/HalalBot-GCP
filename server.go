package HalalBot_GCP

import (
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func init() {
	var err error

	bot, err = linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)

	if err != nil {
		log.Fatal(err)
	}
}

func HalalBot(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Println(err)
				}
			case *linebot.ImageMessage:
				img, err := bot.GetMessageContent(message.ID).Do()
				if err != nil {
					log.Println(err)
				}
				defer img.Content.Close()
				if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(message.OriginalContentURL, message.PreviewImageURL)).Do(); err != nil {
					log.Println(err)
				}
			}
		}
	}
}
