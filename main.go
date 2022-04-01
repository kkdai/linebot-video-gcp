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
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"

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

			// Handle only image message
			case *linebot.ImageMessage:
				log.Println("Got img msg ID:", message.ID)

				//Get image binary from LINE server based on message ID.
				content, err := bot.GetMessageContent(message.ID).Do()
				if err != nil {
					log.Println("Got GetMessageContent err:", err)
				}
				defer content.Content.Close()

				client, err := storage.NewClient(context.Background())
				var ret string
				if err != nil {
					ret = "storage.NewClient: " + err.Error()
				} else {
					ret = "storage.NewClient: OK"
				}

				fileN := buildFileName() + ".jpg"
				if content.ContentLength > 0 {
					uploader := &ClientUploader{
						cl:         client,
						bucketName: bucketName,
						projectID:  projectID,
						uploadPath: "test-files/",
					}

					err = uploader.UploadImage(content.Content)
					if err != nil {
						ret = "uploader.UploadFile: " + err.Error()
					} else {
						ret = "uploader.UploadFile: OK" + " fileN: " + fileN
					}

					imgurl := uploader.GetPulicAddress()

					if _, err = bot.ReplyMessage(event.ReplyToken,
						linebot.NewTextMessage(ret),
						linebot.NewFlexMessage("image",
							&linebot.BubbleContainer{
								Type: linebot.FlexContainerTypeBubble,
								Hero: &linebot.ImageComponent{
									Type: linebot.FlexComponentTypeVideo,
									URL:  imgurl,
								},
								Body: &linebot.BoxComponent{
									Type:   linebot.FlexComponentTypeBox,
									Layout: linebot.FlexBoxLayoutTypeVertical,
									Contents: []linebot.FlexComponent{
										&linebot.TextComponent{
											Type: linebot.FlexComponentTypeText,
											Text: "hello",
										},
									},
								},
							})).Do(); err != nil {
						log.Print(err)
					}

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(ret), linebot.NewImageMessage(imgurl, imgurl)).Do(); err != nil {
						log.Print(err)
					}
				} else {
					log.Println("Empty img")
					ret = "Empty img"

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(ret)).Do(); err != nil {
						log.Print(err)
					}
				}

			// Handle only video message
			case *linebot.VideoMessage:
				log.Println("Got video msg ID:", message.ID)
				ret := "影片上傳判斷中，請稍候"

				go uploadAndDectect(event.Source.UserID, message, bot)

				if _, err = bot.ReplyMessage(event.ReplyToken,
					linebot.NewTextMessage(ret)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func uploadAndDectect(uid string, msg *linebot.VideoMessage, bot *linebot.Client) {
	//Get image binary from LINE server based on message ID.
	content, err := bot.GetMessageContent(msg.ID).Do()
	if err != nil {
		log.Println("Got GetMessageContent err:", err)
	}
	defer content.Content.Close()

	client, err := storage.NewClient(context.Background())
	var ret string
	if err != nil {
		ret = "storage.NewClient: " + err.Error()
	} else {
		ret = "storage.NewClient: OK"
	}

	if content.ContentLength > 0 {
		uploader := &ClientUploader{
			cl:         client,
			bucketName: bucketName,
			projectID:  projectID,
			uploadPath: "test-files/",
		}

		err = uploader.UploadVideo(content.Content)
		if err != nil {
			ret = "uploader.UploadFile: " + err.Error()
		} else {
			ret = "uploader.UploadFile: OK, " + uploader.GetPulicAddress()
		}

		vdourl := uploader.GetPulicAddress()
		flx := newVideoFlexMsg(vdourl, "video link:"+vdourl+" ret:"+ret)

		if _, err = bot.PushMessage(uid, linebot.NewFlexMessage("flex", flx)).Do(); err != nil {
			log.Print(err)
		}
	}
}

func newVideoFlexMsg(video, text string) linebot.FlexContainer {
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Hero: &linebot.VideoComponent{
			Type:       linebot.FlexComponentTypeVideo,
			URL:        video,
			PreviewURL: "https://scdn.line-apps.com/n/channel_devcenter/img/fx/01_1_cafe.png",
			AltContent: &linebot.ImageComponent{
				Type:        linebot.FlexComponentTypeImage,
				URL:         "https://scdn.line-apps.com/n/channel_devcenter/img/fx/01_1_cafe.png",
				Size:        linebot.FlexImageSizeTypeFull,
				AspectRatio: linebot.FlexImageAspectRatioType20to13,
				AspectMode:  linebot.FlexImageAspectModeTypeCover,
			},
			Action: &linebot.URIAction{
				Label: "More information",
				URI:   "http://linecorp.com/",
			},
			AspectRatio: linebot.FlexVideoAspectRatioType20to13,
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: text,
				},
			},
		},
	}
}
