package models

type Message struct {
	From    string `json:"from"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}
