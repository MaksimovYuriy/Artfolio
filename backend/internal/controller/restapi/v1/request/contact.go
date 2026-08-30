package request

import "github.com/maksimovyuriy/artfolio/backend/internal/entity"

type ContactMessage struct {
	SenderEmail string `json:"senderEmail"`
	Message     string `json:"message"`
}

func (r ContactMessage) Entity() entity.ContactMessage {
	return entity.ContactMessage{SenderEmail: r.SenderEmail, Message: r.Message}
}
