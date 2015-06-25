package tgbotapi

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"
)

type UpdateConfig struct {
	Offset  int
	Limit   int
	Timeout int
}

func (bot *Bot) GetUpdates(config UpdateConfig) (chan Update, error) {

	var offset, limit, timeout int

	if config.Offset > 0 {
		offset = config.Offset
	}
	if config.Limit > 0 {
		limit = config.Limit
	}
	if config.Timeout > 0 {
		timeout = config.Timeout
	}

	bot.updates = make(chan Update, 100)

	go func() {
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

				if bot.debug {
					log.Printf("getUpdates: %+v\n", updates)
				}

				for _, e := range updates {
					if e.UpdateId > offset {
						offset = e.UpdateId
					}
					bot.updates <- e
				}
				offset++
			}
		}
	}()

	return bot.updates, nil
}
