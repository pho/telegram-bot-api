package tgbotapi

func NewMessage(chatId int, text string) MessageConfig {
	return MessageConfig{
		ChatId: chatId,
		Text:   text,
		DisableWebPagePreview: false,
		ReplyToMessageId:      0,
	}
}

func NewForward(chatId int, fromChatId int, messageId int) ForwardConfig {
	return ForwardConfig{
		ChatId:     chatId,
		FromChatId: fromChatId,
		MessageId:  messageId,
	}
}

func NewLocation(chatId int, latitude float64, longitude float64) LocationConfig {
	return LocationConfig{
		ChatId:           chatId,
		Latitude:         latitude,
		Longitude:        longitude,
		ReplyToMessageId: 0,
		ReplyMarkup:      nil,
	}
}

func NewChatAction(chatId int, action string) ChatActionConfig {
	return ChatActionConfig{
		ChatId: chatId,
		Action: action,
	}
}

func NewUserProfilePhotos(userId int) UserProfilePhotosConfig {
	return UserProfilePhotosConfig{
		UserId: userId,
		Offset: 0,
		Limit:  0,
	}
}

func NewUpdate(offset int) UpdateConfig {
	return UpdateConfig{
		Offset:  offset,
		Limit:   0,
		Timeout: 0,
	}
}
