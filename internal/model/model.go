package model

import "strings"

type (
	Message struct {
		Sender   string `json:"sender"`
		Receiver string `json:"receiver"`
		Message  string `json:"message"`
	}
)

// IsValid validate the requested message valid for producing
func (m *Message) IsValid() bool {
	return strings.Trim(m.Sender, " ") != "" && strings.Trim(m.Receiver, " ") != "" && strings.Trim(m.Message, " ") != ""
}
