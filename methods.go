package tgbotapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const (
	CHAT_TYPING          = "typing"
	CHAT_UPLOAD_PHOTO    = "upload_photo"
	CHAT_RECORD_VIDEO    = "record_video"
	CHAT_UPLOAD_VIDEO    = "upload_video"
	CHAT_RECORD_AUDIO    = "record_audio"
	CHAT_UPLOAD_AUDIO    = "upload_audio"
	CHAT_UPLOAD_DOCUMENT = "upload_document"
	CHAT_FIND_LOCATION   = "find_location"
)

type MessageConfig struct {
	ChatId                int
	Text                  string
	DisableWebPagePreview bool
	ReplyToMessageId      int
}

type ForwardConfig struct {
	ChatId     int
	FromChatId int
	MessageId  int
}

type LocationConfig struct {
	ChatId           int
	Latitude         float64
	Longitude        float64
	ReplyToMessageId int
	ReplyMarkup      interface{}
}

type ChatActionConfig struct {
	ChatId int
	Action string
}

type UserProfilePhotosConfig struct {
	UserId int
	Offset int
	Limit  int
}

func NewBot(token string, debug bool) *Bot {
	return &Bot{token: token, debug: debug}
}

func (bot *Bot) MakeRequest(endpoint string, params url.Values) (ApiResponse, error) {
	resp, err := http.PostForm("https://api.telegram.org/bot"+bot.token+"/"+endpoint, params)
	defer resp.Body.Close()
	if err != nil {
		return ApiResponse{}, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ApiResponse{}, nil
	}

	if bot.debug {
		log.Println(string(bytes[:]))
	}

	var apiResp ApiResponse
	json.Unmarshal(bytes, &apiResp)

	return apiResp, nil
}

func (bot *Bot) GetMe() (User, error) {
	resp, err := bot.MakeRequest("getMe", nil)
	if err != nil {
		return User{}, err
	}

	var user User
	json.Unmarshal(resp.Result, &user)

	if bot.debug {
		log.Printf("getMe: %+v\n", user)
	}

	bot.self = &user

	return user, nil
}

func (bot *Bot) SendMessage(config MessageConfig) (Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(config.ChatId))
	v.Add("text", config.Text)
	v.Add("disable_web_page_preview", strconv.FormatBool(config.DisableWebPagePreview))
	if config.ReplyToMessageId != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageId))
	}

	resp, err := bot.MakeRequest("sendMessage", v)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.debug {
		log.Printf("sendMessage req : %+v\n", v)
		log.Printf("sendMessage resp: %+v\n", message)
	}

	return message, nil
}

func (bot *Bot) ForwardMessage(config ForwardConfig) (Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(config.ChatId))
	v.Add("from_chat_id", strconv.Itoa(config.FromChatId))
	v.Add("message_id", strconv.Itoa(config.MessageId))

	resp, err := bot.MakeRequest("forwardMessage", v)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.debug {
		log.Printf("forwardMessage req : %+v\n", v)
		log.Printf("forwardMessage resp: %+v\n", message)
	}

	return message, nil
}

func (bot *Bot) SendLocation(config LocationConfig) (Message, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(config.ChatId))
	v.Add("latitude", strconv.FormatFloat(config.Latitude, 'f', 6, 64))
	v.Add("longitude", strconv.FormatFloat(config.Longitude, 'f', 6, 64))
	if config.ReplyToMessageId != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(config.ReplyToMessageId))
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return Message{}, err
		}

		v.Add("reply_markup", string(data))
	}

	resp, err := bot.MakeRequest("sendLocation", v)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.debug {
		log.Printf("sendLocation req : %+v\n", v)
		log.Printf("sendLocation resp: %+v\n", message)
	}

	return message, nil
}

func (bot *Bot) SendChatAction(config ChatActionConfig) error {
	v := url.Values{}
	v.Add("chat_id", strconv.Itoa(config.ChatId))
	v.Add("action", config.Action)

	_, err := bot.MakeRequest("sendChatAction", v)
	if err != nil {
		return err
	}

	return nil
}

func (bot *Bot) GetUserProfilePhotos(config UserProfilePhotosConfig) (UserProfilePhotos, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(config.UserId))
	if config.Offset != 0 {
		v.Add("offset", strconv.Itoa(config.Offset))
	}
	if config.Limit != 0 {
		v.Add("limit", strconv.Itoa(config.Limit))
	}

	resp, err := bot.MakeRequest("getUserProfilePhotos", v)
	if err != nil {
		return UserProfilePhotos{}, err
	}

	var profilePhotos UserProfilePhotos
	json.Unmarshal(resp.Result, &profilePhotos)

	if bot.debug {
		log.Printf("getUserProfilePhotos req : %+v\n", v)
		log.Printf("getUserProfilePhotos resp: %+v\n", profilePhotos)
	}

	return profilePhotos, nil
}

func (bot *Bot) SetWebhook(v url.Values) error {
	_, err := bot.MakeRequest("setWebhook", v)

	return err
}

func (bot *Bot) ClearWebhook() error {
	_, err := bot.MakeRequest("setWebhook", url.Values{})

	return err
}
