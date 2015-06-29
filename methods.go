package tgbotapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

type AudioConfig struct {
	ChatId           int
	ReplyToMessageId int
	ReplyMarkup      interface{}

	UseExistingAudio bool
	FilePath         string
	FileId           string
}

type PhotoConfig struct {
	ChatId           int
	Caption          string
	ReplyToMessageId int
	ReplyMarkup      interface{}

	UseExistingPhoto bool
	FilePath         string
	FileId           string
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
		return ApiResponse{}, err
	}

	if bot.debug {
		log.Println(string(bytes[:]))
	}

	var apiResp ApiResponse
	json.Unmarshal(bytes, &apiResp)

	if !apiResp.Ok {
		return apiResp, errors.New(apiResp.Description)
	}

	return apiResp, nil
}

func (bot *Bot) UploadFile(endpoint string, params map[string]string, fieldname string, filename string) (ApiResponse, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	pwd, err := os.Getwd()
	if err != nil {
		return ApiResponse{}, err
	}

	fmt.Println("Uploading file:", filepath.Base(filename), filepath.FromSlash(pwd+"/"+filename))
	f, err := os.Open(filepath.FromSlash(pwd + "/" + filename))
	if err != nil {
		return ApiResponse{}, err
	}

	fw, err := w.CreateFormFile(fieldname, filepath.Base(filename))
	if err != nil {
		return ApiResponse{}, err
	}

	if _, err = io.Copy(fw, f); err != nil {
		return ApiResponse{}, err
	}

	for key, val := range params {
		if fw, err = w.CreateFormField(key); err != nil {
			return ApiResponse{}, err
		}

		if _, err = fw.Write([]byte(val)); err != nil {
			return ApiResponse{}, err
		}
	}

	w.Close()

	fmt.Println(fw, b)

	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+bot.token+"/"+endpoint, &b)
	if err != nil {
		return ApiResponse{}, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return ApiResponse{}, err
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ApiResponse{}, err
	}

	if bot.debug {
		log.Println(string(bytes[:]))
	}

	var apiResp ApiResponse
	json.Unmarshal(bytes, &apiResp)

	fmt.Println("Upload completed", apiResp)

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

func (bot *Bot) SendAudio(config AudioConfig) (Message, error) {
	if config.UseExistingAudio {
		v := url.Values{}
		v.Add("chat_id", strconv.Itoa(config.ChatId))
		v.Add("audio", config.FileId)
		if config.ReplyToMessageId != 0 {
			v.Add("reply_to_message_id", strconv.Itoa(config.ChatId))
		}
		if config.ReplyMarkup != nil {
			data, err := json.Marshal(config.ReplyMarkup)
			if err != nil {
				return Message{}, err
			}

			v.Add("reply_markup", string(data))
		}

		resp, err := bot.MakeRequest("sendAudio", v)
		if err != nil {
			return Message{}, err
		}

		var message Message
		json.Unmarshal(resp.Result, &message)

		if bot.debug {
			log.Printf("sendAudio req : %+v\n", v)
			log.Printf("sendAudio resp: %+v\n", message)
		}

		return message, nil
	}

	params := make(map[string]string)
	params["chat_id"] = strconv.Itoa(config.ChatId)

	if config.ReplyToMessageId != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageId)
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return Message{}, err
		}

		params["reply_markup"] = string(data)
	}

	fmt.Println("Params:", params)
	resp, err := bot.UploadFile("sendAudio", params, "audio", config.FilePath)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.debug {
		log.Printf("sendAudio resp: %+v\n", message)
	}

	return message, nil
}

func (bot *Bot) SendPhoto(config PhotoConfig) (Message, error) {
	if config.UseExistingPhoto {
		v := url.Values{}
		v.Add("chat_id", strconv.Itoa(config.ChatId))
		v.Add("photo", config.FileId)
		if config.Caption != "" {
			v.Add("caption", config.Caption)
		}
		if config.ReplyToMessageId != 0 {
			v.Add("reply_to_message_id", strconv.Itoa(config.ChatId))
		}
		if config.ReplyMarkup != nil {
			data, err := json.Marshal(config.ReplyMarkup)
			if err != nil {
				return Message{}, err
			}

			v.Add("reply_markup", string(data))
		}

		resp, err := bot.MakeRequest("sendPhoto", v)
		if err != nil {
			return Message{}, err
		}

		var message Message
		json.Unmarshal(resp.Result, &message)

		if bot.debug {
			log.Printf("sendPhoto req : %+v\n", v)
			log.Printf("sendPhoto resp: %+v\n", message)
		}

		return message, nil
	}

	params := make(map[string]string)
	params["chat_id"] = strconv.Itoa(config.ChatId)
	if config.Caption != "" {
		params["caption"] = config.Caption
	}
	if config.ReplyToMessageId != 0 {
		params["reply_to_message_id"] = strconv.Itoa(config.ReplyToMessageId)
	}
	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return Message{}, err
		}

		params["reply_markup"] = string(data)
	}

	resp, err := bot.UploadFile("sendPhoto", params, "photo", config.FilePath)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.debug {
		log.Printf("sendPhoto resp: %+v\n", message)
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
