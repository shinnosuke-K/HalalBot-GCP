package HalalBot_GCP

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"

	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	bot       *linebot.Client
	hl        *halalFood
	lineStamp map[bool]map[string]string
)

func init() {
	var err error

	bot, err = linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)

	if err != nil {
		log.Fatal(err)
	}

	hl = newHalal()

	lineStamp = map[bool]map[string]string{
		true: {
			"packageID": "2",
			"stickerID": "179",
		},
		false: {
			"packageID": "2",
			"stickerID": "39",
		},
	}
}

func ocr(ctn io.ReadCloser) ([]string, bool) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	img, err := vision.NewImageFromReader(ctn)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	texts, err := client.DetectTexts(ctx, img, nil, 1)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	var detectedTextLists []string
	for _, text := range texts {
		detectedTextLists = append(detectedTextLists, text.GetDescription())
	}

	return detectedTextLists, true
}

type halalFood struct {
	ngFoods []string
}

func newHalal() *halalFood {
	return &halalFood{ngFoods: []string{"ワイン", "みりん", "日本酒", "ビール", "ラム酒", "料理酒", "豚肉", "豚", "ポーク", "ゼラチン", "ラード"}}
}

func (hf *halalFood) judge(texts []string) (string, bool) {
	for _, text := range texts {
		log.Println(text)
		if ok := hf.in(text); ok {
			return text, false
		}
	}
	return "", true
}

func (hf *halalFood) in(word string) bool {
	for _, food := range hf.ngFoods {
		if ok := strings.Contains(word, food); ok {
			return true
		}
	}
	return false
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
				texts, ok := ocr(ctn.Content)

				switch ok {
				case true:
					foodName, canEat := hl.judge(texts)
					if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewStickerMessage(lineStamp[canEat]["packageID"], lineStamp[canEat]["stickerID"])).Do(); err != nil {
						log.Println(err)
					}

					if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(foodName)).Do(); err != nil {
						log.Println(err)
					}

				case false:
					if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Please retry")).Do(); err != nil {
						log.Println(err)
					}
				}

			}
		}
	}
}
