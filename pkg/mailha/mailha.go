package mailha

import (
	"gopkg.in/gomail.v2"
)

type MessageWrapper struct {
	CC     []string
	To     []string
	From   string
	Title  string
	Body   string
	Attach []string
}

type IMailSender interface {
	SendMail(MessageWrapper) error
}

type IMessageBuilder interface {
	SetTitle(string)
	SetBody(string)
	SetFromMail(string)
	AddToMail(...string)
	AddCC(...string)
	AddAttach(...string)
	Build() (MessageWrapper, error)
}
type MailSenderConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
}
type MailSender struct {
	config MailSenderConfig
}

func NewMailSender(config MailSenderConfig) *MailSender {
	return &MailSender{config: config}
}

type MessageBuilder struct {
	message MessageWrapper
}

func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{message: MessageWrapper{}}
}

func (mw *MessageWrapper) convert() *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", mw.From)
	m.SetHeader("To", mw.To...)
	if len(mw.CC) > 0 {
		m.SetHeader("Cc", mw.CC...)
	}
	m.SetHeader("Subject", mw.Title)
	m.SetBody("text/html", mw.Body)
	for _, attach := range mw.Attach {
		m.Attach(attach)
	}
	return m
}

func (mb *MessageBuilder) SetTitle(title string) {
	mb.message.Title = title
}

func (mb *MessageBuilder) SetBody(body string) {
	mb.message.Body = body
}

func (mb *MessageBuilder) SetFromMail(from string) {
	mb.message.From = from
}

func (mb *MessageBuilder) AddToMail(to ...string) {
	mb.message.To = append(mb.message.To, to...)
}

func (mb *MessageBuilder) AddCC(cc ...string) {
	mb.message.CC = append(mb.message.CC, cc...)
}

func (mb *MessageBuilder) AddAttach(attach ...string) {
	mb.message.Attach = append(mb.message.Attach, attach...)
}

func (mb *MessageBuilder) Build() (MessageWrapper, error) {
	// Here you can add any additional validation or processing if needed
	return mb.message, nil
}
