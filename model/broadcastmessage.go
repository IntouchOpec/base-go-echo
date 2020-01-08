package model

import "time"

type MessageType string

const (
	MessageTypeMassage MessageType = "Massage"
	MessageTypeImage   MessageType = "Image"
	MessageTypeVideo   MessageType = "Video"
	MessageTypeAudio   MessageType = "Audio"
	MessageTypeFlex    MessageType = "Flex"
)

type SendState int

const (
	ByOne          SendState = 1
	ByCustomerGrop SendState = 2
	All            SendState = 3
	Tester         SendState = 4
)

type SandDateType int

const (
	Now     SandDateType = 1
	SetTime SandDateType = 2
)

type BroadcastMessage struct {
	SandDate      time.Time    `json:"time_sand"`
	SandDateType  SandDateType `json:"sand_date_type"`
	Send          string       `json:"send"`
	SendState     SendState    `json:"send_state"`
	Message       string       `json:"message"`
	MessageType   MessageType  `json:"message_type" gorm:"type:varchar(25)"`
	AccountID     uint         `json:"account_id"`
	Account       *Account     `json:"account" gorm:"ForeignKey:AccountID"`
	ChatChannelID uint         `json:"chat_channel_id"`
	ChatChannel   ChatChannel  `json:"chat_channel" gorm:"ForeignKey:ChatChannelID"`
}
