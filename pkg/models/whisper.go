package models

type Whisper struct {
	Message
	To string `json:"to"`
}
