package web

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	guuid "github.com/google/uuid"
	"github.com/labstack/echo"
)

// serviceListHandler
func ServiceListHandler(c *Context) error {
	services := []*model.Service{}
	a := auth.Default(c)
	queryPar := c.QueryParams()
	page, limit := SetPagination(queryPar)
	var total int
	db := model.DB()

	filterService := db.Where("ser_account_id = ?", a.User.GetAccountID()).Find(&services).Count(&total)

	filterService.Limit(limit).Offset(page).Find(&services)
	pagination := MakePagination(total, page, limit)

	err := c.Render(http.StatusOK, "service-list", echo.Map{
		"list":       services,
		"title":      "service",
		"pagination": pagination,
	})
	return err
}

// serviceDetailHandler
func ServiceDetailHandler(c *Context) error {
	service := model.Service{}
	id := c.Param("id")
	a := auth.Default(c)
	model.DB().Preload("Account").Where("ser_account_id = ? ", a.User.GetAccountID()).Find(&service, id)
	err := c.Render(http.StatusOK, "service-detail", echo.Map{
		"detail": service,
		"title":  "service",
	})
	return err
}

func ServiceCreateHandler(c *Context) error {
	Service := model.Service{}
	csrfValue := c.Get("_csrf")

	err := c.Render(http.StatusOK, "service-form", echo.Map{
		"detail": Service,
		"title":  "service",
		"_csrf":  csrfValue,
	})
	return err
}

func ServiceEditHandler(c *Context) error {
	service := model.Service{}
	id := c.Param("id")
	a := auth.Default(c)
	model.DB().Preload("Account").Preload("ServiceSlots").Preload("ChatChannels").Where("ser_account_id = ? ", a.User.GetAccountID()).Find(&service, id)
	err := c.Render(http.StatusOK, "service-form", echo.Map{
		"detail": service,
		"title":  "service",
	})
	return err
}

func ServiceDeleteHandler(c *Context) error {
	id := c.Param("id")
	pro := model.DeleteserviceByID(id)
	err := c.JSON(http.StatusOK, pro)
	return err
}

// func ServiceSlotCreateHandler(c *Context) error {
// 	messageTypes := []linebot.MessageType{linebot.MessageTypeText, linebot.MessageTypeImage, linebot.MessageTypeVideo, linebot.MessageTypeAudio, linebot.MessageTypeFile, linebot.MessageTypeLocation, linebot.MessageTypeSticker, linebot.MessageTypeTemplate, linebot.MessageTypeImagemap, linebot.MessageTypeFlex}

// 	sunService := model.ServiceSlot{}
// 	err := c.Render(http.StatusOK, "sub-service-form", echo.Map{
// 		"detail":       sunService,
// 		"title":        "service",
// 		"messageTypes": messageTypes,
// 	})
// 	return err
// }

type serviceForm struct {
	Name   string  `form:"name"`
	Detail string  `form:"detail"`
	Price  float32 `form:"price"`
	// Image  byte   `form:"file"`
}

var (
	ErrBucket       = errors.New("Invalid bucket!")
	ErrSize         = errors.New("Invalid size!")
	ErrInvalidImage = errors.New("Invalid image!")
)

func ServicePostHandler(c *Context) error {
	service := serviceForm{}
	file := c.FormValue("file")

	idx := strings.Index(file, ";base64,")
	if idx < 0 {
		// return "", ErrInvalidImage
	}
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(file[idx+8:]))
	buff := bytes.Buffer{}
	_, err := buff.ReadFrom(reader)
	if err != nil {
		// return "", err
	}
	imgCfg, fm, err := image.DecodeConfig(bytes.NewReader(buff.Bytes()))
	if err != nil {
		// return "", err
	}

	if imgCfg.Width != 750 || imgCfg.Height != 685 {
		// return "", ErrSize
	}
	if fm == "" {
		fm = ".jpg"
	}

	u := guuid.New()
	fileNameBase := "public/assets/images/%s"
	fileNameBase = fmt.Sprintf(fileNameBase, u)
	fileName := fileNameBase + fm
	err = ioutil.WriteFile(fileName, buff.Bytes(), 0644)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Bind(&service); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	a := auth.Default(c)
	serviceModel := model.Service{
		SerName:   service.Name,
		SerDetail: service.Detail,
		SerPrice:  service.Price,
		SerImage:  fmt.Sprintf("%s.%s", u, fm),
		AccountID: a.User.GetAccountID(),
	}
	serviceModel.Saveservice()

	return c.JSON(http.StatusCreated, serviceModel)
}

type ServiceSlotForm struct {
	Start  string `form:"start" json:"start"`
	End    string `form:"end" json:"end"`
	Day    int    `form:"day" json:"day"`
	Amount int    `form:"amount" json:"amount"`
}

// func ServiceSlotPostHandler(c *Context) error {
// 	id := c.Param("id")
// 	Service := model.Service{}
// 	db := model.DB()
// 	serviceSlotFrom := ServiceSlotForm{}
// 	if err := c.Bind(&serviceSlotFrom); err != nil {
// 		return c.JSON(http.StatusBadRequest, err)
// 	}
// 	db.Find(&Service, id)
// 	db.Model(&Service).Association("ServiceSlots").Append(&model.ServiceSlot{Start: serviceSlotFrom.Start, End: serviceSlotFrom.End, Day: serviceSlotFrom.Day, Amount: serviceSlotFrom.Amount})
// 	return c.JSON(http.StatusCreated, Service)
// }

// func ServiceSlotEditHandler(c *Context) error {
// 	sunService := model.ServiceSlot{}
// 	id := c.Param("id")
// 	a := auth.Default(c)
// 	messageTypes := []linebot.MessageType{linebot.MessageTypeText, linebot.MessageTypeImage, linebot.MessageTypeVideo, linebot.MessageTypeAudio, linebot.MessageTypeFile, linebot.MessageTypeLocation, linebot.MessageTypeSticker, linebot.MessageTypeTemplate, linebot.MessageTypeImagemap, linebot.MessageTypeFlex}
// 	model.DB().Preload("service", func(db *gorm.DB) *gorm.DB {
// 		return db.Preload("ChatChannels").Where("account_id = ? ", a.User.GetAccountID())
// 	}).Find(&sunService, id)
// 	err := c.Render(http.StatusOK, "sub-service-form", echo.Map{
// 		"detail":       sunService,
// 		"title":        "service",
// 		"messageTypes": messageTypes,
// 	})
// 	return err
// }

// func ServiceSlotDeleteHandler(c *Context) error {
// 	id := c.Param("id")
// 	serviceSlot := model.DeleteServiceSlot(id)
// 	return c.JSON(http.StatusOK, serviceSlot)
// }

func ServiceChatChannelViewHandler(c *Context) error {
	chatChannels := []*model.ChatChannel{}
	a := auth.Default(c)
	model.DB().Preload("Account").Where("account_id = ?", a.User.GetAccountID()).Find(&chatChannels)

	err := c.Render(http.StatusOK, "service-chat-channel-form", echo.Map{
		"list_chat_channel": chatChannels,
		"title":             "service",
	})
	return err
}

func ServiceChatChannelPostHandler(c *Context) error {
	id := c.QueryParam("id")
	chatChannelID := c.FormValue("chat_channel_id")
	service := model.Service{}
	chatChannel := model.ChatChannel{}
	db := model.DB()

	if err := db.Find(&chatChannel, chatChannelID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := db.Find(&service, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := db.Model(&service).Association("ChatChannels").Append(chatChannel).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, service)
}
