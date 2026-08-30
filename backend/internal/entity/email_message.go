package entity

type EmailMessage struct {
	Recipient string `json:"recipient"`
	ReplyTo   string `json:"replyTo"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}
