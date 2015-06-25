package tgbotapi

import (
	"fmt"
	"testing"
)

func TestEchobot(t *testing.T) {

	//Create a new bot
	b := NewBot("YourTokenHere", false)

	//Get the updates channel
	c, _ := b.GetUpdates(UpdateConfig{Timeout: 60})

	fmt.Println("Waiting for updates...")
	if e, ok := <-c; ok {
		fmt.Println("Someone said:", e.Message.Text)

		b.SendChatAction(NewChatAction(e.Message.From.Id, CHAT_TYPING))
		b.SendMessage(MessageConfig{ChatId: e.Message.From.Id, Text: e.Message.Text})
	}
	fmt.Println("Bye")

}
