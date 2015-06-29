package tgbotapi

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"
)

// GetUpdates returns a chan filled whenever a new update is gotten.
func (bot *Bot) GetUpdatesChan(config UpdateConfig) (chan Update, error) {

	offset := config.Offset
	limit := config.Limit
	timeout := config.Timeout

	bot.updates = make(chan Update, 100)

	go func() {
		defer close(bot.updates)

		for {
			v := url.Values{}
			if offset > 0 {
				v.Add("offset", strconv.Itoa(offset))
			}
			if limit > 0 {
				v.Add("limit", strconv.Itoa(limit))
			}
			if timeout > 0 {
				v.Add("timeout", strconv.Itoa(timeout))
			}

			resp, err := bot.MakeRequest("getUpdates", v)
			if err == nil {
				var updates []Update
				json.Unmarshal(resp.Result, &updates)

				if bot.Debug {
					log.Printf("getUpdates: %+v\n", updates)
				}

				for _, e := range updates {
					if e.UpdateID >= offset {
						offset = e.UpdateID + 1
					}

					bot.updates <- e
				}
			}
		}
	}()

	return bot.updates, nil

}
