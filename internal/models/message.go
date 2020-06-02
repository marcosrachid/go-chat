package models

type Message struct {
	Channel string `json:"channel"`
	From    string `json:"from"`
	Text    string `json:"text"`
}
