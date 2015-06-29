package tgbotapi

import (
	"fmt"
	"testing"
)

func TestEchobot(t *testing.T) {

	//Create a new bot
	if b, err := NewBot("YourAwesomebotToken"); err == nil {

		//Get the updates channel
		if c, err := b.GetUpdatesChan(UpdateConfig{Timeout: 60}); err == nil {

			fmt.Println("Waiting for updates...")
			if e, ok := <-c; ok {
				fmt.Println("Someone said:", e.Message.Text)

				// Reply with the same message
				if _, err := b.SendMessage(MessageConfig{ChatID: e.Message.Chat.ID, Text: e.Message.Text}); err != nil {
					t.Error("Failed sending the message")
				}

			} else {
				t.Error("Failed getting any updates")
			}

		} else {
			t.Error("Failed getting updates chan")
		}

	} else {
		t.Error("Failed creating the bot")
	}
}
