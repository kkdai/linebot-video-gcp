LINE Bot Video with Google Cloud Storage and STT: A LINE Bot using GCS and Speech-To-Text to translate your video with video script.
==============

 [![GoDoc](https://godoc.org/github.com/kkdai/linebot-video-gcp.svg?status.svg)](https://godoc.org/github.com/kkdai/LineBotTemplate)   ![Go](https://github.com/kkdai/linebot-video-gcp/workflows/Go/badge.svg) [![goreportcard.com](https://goreportcard.com/badge/github.com/kkdai/linebot-video-gcp)](https://goreportcard.com/report/github.com/kkdai/linebot-video-gcp)

# 解決的問題痛點

不知道大家是否有類似的問題，有一些長輩或是朋友們由於可能不方便打字，經常都透過語音訊息來溝通。但是許多時候其實是不方便打開手機音量來「聽訊息」或者是「看語音影片」的，比如說在捷運上，或是在會議之中。 這個時候，經常會想到說是否有一個方便的服務可以透過上傳影片（或是語音檔案）的方式。 透過 STT ( Speech-To-Text )的技術來幫影片中的語音翻譯出文字，並且透過 LINE 這次開放的 [Flex Message Update 3](https://developers.line.biz/en/news/2022/03/11/flex-message-update-3-released/) 新格式 ，來讓使用者更快看得到你的訊息內容，而不需要打開影片內容。

# 如何用？如何架設？ :


## 開源程式碼：
不囉唆，先來看 code 。

#### 開源程式碼： [https://github.com/kkdai/linebot-video-gcp]( https://github.com/kkdai/linebot-video-gcp)

## 如何快速架設在自己平台：

請根據以下的流程，來直接透過開源的程式碼來架設一個 LINE Bot 跟 GCP  STT 與 GCS (Google Cloud Storage) 的結合。

### 事前準備

- [LINE Developers 帳號](https://developers.line.biz/console/)
  - 到 [LINE Developers Console](https://developers.line.biz/console/) 透過你的 LINE 帳號登入。
  - 開啟 LINE Developers 帳號
  - 開啟一個官方帳號 Message API Channel 並且取得 `channel secret` 與 `channel access token`。 請[參考這篇教學](https://developers.line.biz/en/docs/messaging-api/getting-started/#using-console)。
- 需要一個免費的 [Heroku](https://dashboard.heroku.com/) 帳號。
- 需要一個收費的 [Google Cloud Platform](https://cloud.google.com/) 帳號，並且取得 GCP JSON key file ，更多詳細流程請參考 [Cloud Storage client libraries](https://cloud.google.com/storage/docs/reference/libraries)。

### 開始部署

1. 請先到程式碼:  [https://github.com/kkdai/linebot-video-gcp]( https://github.com/kkdai/linebot-video-gcp)
2. 按下： [![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)  來部署到你的 Heroku 帳號
   
3. 除了填寫 App Name 以外，以下的參數必須要填入才能完整運行。
   ![image-20220421101405675](http://www.evanlin.com/images/2021/image-20220421101405675.png)
   
   1. **ChannelAccessToken**: 請到 LINE Developers Console issue 一個。
   2. **ChannelSecret**: 請到 LINE Developers Console 拿一個。
   3. **GCS_BUCKET_NAME**: Google Cloud Storage Bucket Name 你需要先建立一個 GCS bucket 並且依照稍後文章的內容來改成 public access 。
   4. **GCS_PROJECT_ID**: Google Cloud Storage Project ID
   5. **GOOGLE_APPLICATION_CREDENTIALS**: 這是必須要填寫，請勿修改。
   6. **GOOGLE_CREDENTIALS**: 這邊請將 GCP JSON key file 的內容（對～就是整包檔案裡面的 JSON 資料），全部放上去。這樣是比較安全的。
4. 請到 LINE 官方帳號的平台，到了右上角的「設定」中，選擇「帳號設定」
   1. 將你官方帳號基本資料設定好，並且打開加入群組功能。
      1. ![image-20220421103018014](http://www.evanlin.com/images/2021/image-20220421103018014.png)
   
   2. 到回應設定，將以下設定改好：
      1. ![image-20220421103154889](http://www.evanlin.com/images/2021/image-20220421103154889.png)
      2. 回應模式改成「聊天機器人」
      3. 停用「自動回應訊息」
      4. 啟用「Webhook」
   3. 到 Messaging API 選項，將 Webhook 網址填入`: https://{YOUR_HEROKU_SERVER_ID}.herokuapp.com/callback
   4. 關於如何快速部署的流程，可以參考另外一篇文章內的影片：
     - [How to deploy LINE BotTemplate](https://www.youtube.com/watch?v=0BIknEz1f8k)
     - [Hoe to modify your LINE BotTemplate code](https://www.youtube.com/watch?v=ckij73sIRik)
   

### 如何使用:

- 打開聊天機器人
- 傳送影片訊息
- 過一兩秒的判別時間，就會收到結果的 Flex Message

### 成果

![image-20220407212250046](http://www.evanlin.com/images/2021/image-20220407212250046.png)



### 完整開發教學：

[[學習文件] 如何透過 GCP STT (Speech To Text) 與 LINE Video Flex Message 來打造幫助長輩的影片翻譯（自動輸出文字）聊天機器人](https://www.evanlin.com/til-gcp-speech-video-flex/)



License
---------------

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

