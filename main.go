// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client

var projectID string
var bucketName string

func main() {
	var err error
	projectID = os.Getenv("GCS_PROJECT_ID")
	bucketName = os.Getenv("GCS_BUCKET_NAME")

	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			// Handle only on text message
			case *linebot.TextMessage:
				ret := message.Text
				if message.Text == "video" {
					_, err := storage.NewClient(context.Background())
					if err != nil {
						ret = "storage.NewClient: " + err.Error()
					} else {
						ret = "storage.NewClient: OK"
					}
				}
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(ret)).Do(); err != nil {
					log.Print(err)
				}
			// Handle only on Sticker message
			case *linebot.StickerMessage:
				var kw string
				for _, k := range message.Keywords {
					kw = kw + "," + k
				}

				outStickerResult := fmt.Sprintf("收到貼圖訊息: %s, pkg: %s kw: %s  text: %s", message.StickerID, message.PackageID, kw, message.Text)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(outStickerResult)).Do(); err != nil {
					log.Print(err)
				}

			// Handle only video message
			case *linebot.VideoMessage:
				ret := fmt.Sprintf("Get video from %s", message.ContentProvider.OriginalContentURL)
				client, err := storage.NewClient(context.Background())
				if err != nil {
					ret = "storage.NewClient: " + err.Error()
				} else {
					ret = "storage.NewClient: OK"
				}

				bodyBytes, _ := io.ReadAll(r.Body)

				log.Println("Got video raw event data from:", string(bodyBytes))
				log.Println("Got video msg from:", message)
				log.Println("Got file from:", message.ContentProvider.OriginalContentURL)

				if len(message.ContentProvider.OriginalContentURL) > 0 {
					// Get the video data
					resp, err := http.Get(message.ContentProvider.OriginalContentURL)
					if err != nil {
						log.Print(err)
					}
					defer resp.Body.Close()

					uploader := &ClientUploader{
						cl:         client,
						bucketName: bucketName,
						projectID:  projectID,
						uploadPath: "test-files/",
					}

					err = uploader.UploadFile(resp.Body, "video.mp4")
					if err != nil {
						ret = "uploader.UploadFile: " + err.Error()
					} else {
						ret = "uploader.UploadFile: OK"
					}

				} else {
					log.Println("Empty video")
					ret = "Empty video"
				}

				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(ret)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
