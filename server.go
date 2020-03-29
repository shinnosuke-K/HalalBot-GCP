package HalalBot_GCP

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	vision "cloud.google.com/go/vision/apiv1"

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

func ocr(ctn io.ReadCloser) string {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Println(err)
		return "Not Annotate Client"
	}

	img, err := vision.NewImageFromReader(ctn)
	if err != nil {
		log.Println(err)
		return "Not Image From Reader"
	}

	texts, err := client.DetectTexts(ctx, img, nil, 1)
	if err != nil {
		log.Println(err)
		return "Not Detect Texts"
	}

	return texts[0].GetDescription()
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
				ctn, err := bot.GetMessageContent(message.ID).Do()
				if err != nil {
					log.Println(err)
				}

				defer ctn.Content.Close()
				if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(ocr(ctn.Content))).Do(); err != nil {
					log.Println(err)
				}
			}
		}
	}
}
