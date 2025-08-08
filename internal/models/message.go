package models

import (
	"time"
)

type Message struct {
	Sender    User      `json:"sender"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func NewMessage(sender User, content string) Message {
	message := Message{
		Sender:    sender,
		Content:   content,
		Timestamp: time.Now(),
	}

	return message
}

func (m Message) String() string {
	return "[" + m.Timestamp.Format("2006-01-02 15:04:05") + "] " + m.Sender.Name + ": " + m.Content + "\n"
}
