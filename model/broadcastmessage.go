package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/lib/uploadgoolgestorage"

	// lib "github.com/IntouchOpec/base-go-echo"
	"github.com/hb-go/echo-web/model"

	// "github.com/IntouchOpec/base-go-echo/lib"
	"github.com/line/line-bot-sdk-go/linebot"
)

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

type ActionType string

const (
	ActionTypeBroacast  ActionType = "Broacast"
	ActionTypePromotion ActionType = "Promotion"
)

type BroadcastMessage struct {
	SandDate      time.Time    `json:"time_sand"`
	SendTo        string       `json:"send_to"`
	SandDateType  SandDateType `json:"sand_date_type"`
	Send          string       `json:"send"`
	Respose       string       `json:"respose"`
	SendState     SendState    `json:"send_state"`
	Message       string       `json:"message"`
	ActionType    ActionType   `json:"action_type" gorm:"type:varchar(25)"`
	MessageType   MessageType  `json:"message_type" gorm:"type:varchar(25)"`
	AccountID     uint         `json:"account_id"`
	Account       *Account     `json:"account" gorm:"ForeignKey:AccountID"`
	ChatChannelID uint         `json:"chat_channel_id"`
	ChatChannel   ChatChannel  `json:"chat_channel" gorm:"ForeignKey:ChatChannelID"`
}

func (br *BroadcastMessage) Create() error {
	var message linebot.SendingMessage
	db := model.DB()
	var err error
	switch os := br.MessageType; os {
	case "Massage":
		message = linebot.NewTextMessage(br.Message)
	case "Image":
		image := linebot.ImageMessage{}
		if err := json.Unmarshal([]byte(br.Message), image); err != nil {
			return err
		}
		ctx := context.Background()
		imagePath, err := uploadgoolgestorage.UploadGoolgeStorage(ctx, image.OriginalContentURL, "images/broadcast/")
		urlFile := fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, imagePath)
		if err != nil {
			return err
		}
		message = linebot.NewImageMessage(urlFile, urlFile)
	case "Video":
		video := linebot.VideoMessage{}
		if err := json.Unmarshal([]byte(br.Message), video); err != nil {
			return err
		}
		ctx := context.Background()
		videoPath, err := uploadgoolgestorage.UploadGoolgeStorage(ctx, video.OriginalContentURL, "video/broadcast/")
		if err != nil {
			return err
		}
		urlFile := fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, videoPath)
		if err != nil {
			return err
		}
		message = linebot.NewVideoMessage(urlFile, urlFile)
	case "Audio":
		audio := linebot.AudioMessage{}
		if err := json.Unmarshal([]byte(br.Message), audio.OriginalContentURL); err != nil {
			return err
		}
		ctx := context.Background()
		audio.OriginalContentURL, err = uploadgoolgestorage.UploadGoolgeStorage(ctx, audio.OriginalContentURL, "audio/broadcast/")
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		audio.OriginalContentURL = fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, audio.OriginalContentURL)
		message = linebot.NewAudioMessage(audio.OriginalContentURL, audio.Duration)
	case "Line_Bot_Designer":
		flex := br.Message
		flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(flex))
		if err != nil {
			return err
		}
		message = linebot.NewFlexMessage("flex message", flexContainer)
	}
	out, err := json.Marshal(message)
	if err != nil {
		return err
	}
	br.Message = string(out)
	if err := db.Create(&br).Error; err != nil {
		return err
	}
	return nil
}

func (br *BroadcastMessage) SendMessage() error {
	var err error
	bot, err := linebot.New(br.ChatChannel.ChaChannelSecret, br.ChatChannel.ChaChannelAccessToken)
	if err != nil {
		return err
	}
	var message linebot.SendingMessage
	var rep *linebot.BasicResponse
	var recipient []string
	switch br.SendState {
	case 1:
		customers := []Customer{}
		lineNames := br.SendTo
		db.Where("account_id and cus_display_name = ?", br.AccountID, lineNames).Find(&customers)
		for _, customer := range customers {
			recipient = append(recipient, customer.CusLineID)
		}
		rep, err = bot.Multicast(recipient, message).Do()
		if err != nil {
			return err
		}
	case 2:
		customers := []Customer{}
		customerTypeID := br.SendTo
		db.Preload("CustomerType", "id = ?", customerTypeID).Find(&customers)
		var recipient []string
		for _, customer := range customers {
			recipient = append(recipient, customer.CusLineID)
		}
		rep, err = bot.Multicast(recipient, message).Do()
		if err != nil {
			return err
		}

	case 3:
		rep, err = bot.BroadcastMessage(message).Do()
		if err != nil {
			return err
		}
	case 4:
		var testers []User
		db.Where("tester = ? and account_id", true, br.ChatChannelID).Find(&testers)
		var recipient []string
		for _, tester := range testers {
			recipient = append(recipient, tester.LineID)
		}
		rep, err = bot.Multicast(recipient, message).Do()
	}
	repB, err := json.Marshal(&rep)
	if err != nil {
		return err
	}
	br.Respose = string(repB)
	if err := db.Update(&br).Error; err != nil {
		return err
	}
	return nil
}
