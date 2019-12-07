package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
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

	filterService := db.Model(&services).Where("account_id = ?", a.User.GetAccountID()).Count(&total)

	pagination := MakePagination(total, page, limit)
	filterService.Limit(pagination.Record).Offset(pagination.Offset).Find(&services)

	err := c.Render(http.StatusOK, "service-list", echo.Map{
		"list":       services,
		"title":      "service",
		"pagination": pagination,
	})
	return err
}

func ServiceEditViewHandler(c *Context) error {
	Service := model.Service{}
	id := c.Param("id")
	accID := auth.Default(c).GetAccountID()

	model.DB().Where("account_id = ?", accID).Find(&Service, id)
	return c.Render(http.StatusOK, "service-form", echo.Map{
		"detail": Service,
		"title":  "service",
		"method": "PUT",
	})
}

func ServiceEditPutHandler(c *Context) error {
	service := serviceForm{}
	image := c.FormValue("image")
	id := c.Param("id")
	var err error
	if image == "" {
		file := c.FormValue("file")
		image, _, err = lib.UploadteImage(file)
	}

	if err := c.Bind(&service); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	a := auth.Default(c)
	serviceModel := model.Service{
		SerName:   service.Name,
		SerDetail: service.Detail,
		SerPrice:  service.Price,
		SerTime:   service.Time,
		SerImage:  image,
		AccountID: a.User.GetAccountID(),
	}
	err = serviceModel.UpdateService(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"data":     serviceModel,
		"redirect": fmt.Sprintf("/admin/service/%d", serviceModel.ID),
	})
}

// serviceDetailHandler
func ServiceDetailHandler(c *Context) error {
	service := model.Service{}
	id := c.Param("id")
	a := auth.Default(c)
	model.DB().Where("account_id = ? ", a.User.GetAccountID()).Find(&service, id)
	err := c.Render(http.StatusOK, "service-detail", echo.Map{
		"detail": service,
		"title":  "service",
	})
	return err
}

func ServiceCreateHandler(c *Context) error {
	Service := model.Service{}

	err := c.Render(http.StatusOK, "service-form", echo.Map{
		"detail": Service,
		"title":  "service",
		"method": "POST",
	})
	return err
}

func ServiceEditHandler(c *Context) error {
	service := model.Service{}
	id := c.Param("id")
	a := auth.Default(c)
	model.DB().Preload("Account").Preload("ServiceSlots").Preload("ChatChannels").Where("account_id = ? ", a.User.GetAccountID()).Find(&service, id)
	err := c.Render(http.StatusOK, "service-form", echo.Map{"method": "PUT",
		"detail": service,
		"title":  "service",
	})
	return err
}

func ServiceDeleteImageHandler(c *Context) error {
	service := model.Service{}
	id := c.Param("id")
	accID := auth.Default(c).GetAccountID()

	if err := model.DB().Where("account_id = ? ", accID).Find(&service, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := lib.DeleteFile(service.SerImage); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := service.RemoveImage(id); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"detail": service,
	})
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
// 	err := c.Render(http.StatusOK, "sub-service-form", echo.Map{"method": "PUT",
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
	Time   string  `form:"time"`
}

var (
	ErrBucket       = errors.New("Invalid bucket!")
	ErrSize         = errors.New("Invalid size!")
	ErrInvalidImage = errors.New("Invalid image!")
)

func ServicePostHandler(c *Context) error {
	service := serviceForm{}
	file := c.FormValue("file")
	fileUrl, _, err := lib.UploadteImage(file)

	if err := c.Bind(&service); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	a := auth.Default(c)
	serviceModel := model.Service{
		SerName:   service.Name,
		SerDetail: service.Detail,
		SerPrice:  service.Price,
		SerTime:   service.Time,
		SerImage:  fileUrl,
		AccountID: a.User.GetAccountID(),
	}
	err = serviceModel.SaveService()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"data":     serviceModel,
		"redirect": fmt.Sprintf("/admin/service/%d", serviceModel.ID),
	})
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
// 	err := c.Render(http.StatusOK, "sub-service-form", echo.Map{"method": "PUT",
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

	err := c.Render(http.StatusOK, "service-chat-channel-form", echo.Map{"method": "PUT",
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
