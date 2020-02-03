package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
		ctx := context.Background()
		image, err = lib.UploadGoolgeStorage(ctx, file, "images/Service/")
	}

	if err := c.Bind(&service); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	accID := auth.Default(c).GetAccountID()
	serviceModel := model.Service{}
	db := model.DB()
	if err := db.Where("account_id = ?", accID).Find(&serviceModel, id).Error; err != nil {
		return err
	}
	serviceModel.SerName = service.Name
	serviceModel.SerDetail = service.Detail
	serviceModel.SerPrice = service.Price
	serviceModel.SerTime = service.Time
	serviceModel.SerImage = image

	if err := db.Save(&serviceModel).Error; err != nil {
		return err
	}

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
	model.DB().Preload("ServiceItems").Where("account_id = ? ", a.User.GetAccountID()).Find(&service, id)
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
	ctx := context.Background()
	if _, err := lib.RemoveFileGoolgeStorage(ctx, "triple-t", service.SerImage); err != nil {
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
	Price  float64 `form:"price"`
	Time   string  `form:"time"`
}

var (
	ErrBucket       = errors.New("Invalid bucket!")
	ErrSize         = errors.New("Invalid size!")
	ErrInvalidImage = errors.New("Invalid image!")
)

func ServicePatchHandler(c *Context) error {
	service := model.Service{}
	id := c.Param("id")

	if err := model.DB().Find(&service, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := c.Bind(&service); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, service)
}

func ServicePostHandler(c *Context) error {
	service := serviceForm{}
	file := c.FormValue("file")
	ctx := context.Background()
	imagePath, err := lib.UploadGoolgeStorage(ctx, file, "images/Service/")

	if err := c.Bind(&service); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	a := auth.Default(c)
	serviceModel := model.Service{
		SerName:   service.Name,
		SerDetail: service.Detail,
		SerPrice:  service.Price,
		SerTime:   service.Time,
		SerImage:  imagePath,
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

func ServiceItemCreateViewHandlder(c *Context) error {
	var service model.Service
	var serviceItem model.ServiceItem
	return c.Render(http.StatusOK, "service-item-form", echo.Map{
		"service": service,
		"detail":  serviceItem,
		"title":   "service",
		"method":  "POST",
	})
}

type ServiceItemReq struct {
	Time  string  `form:"time"`
	Price float64 `form:"price"`
	Name  string  `form:"name"`
}

func ServiceItemCreatePostHandler(c *Context) error {
	id := c.Param("id")
	accID := auth.Default(c).GetAccountID()
	db := model.DB()
	var service model.Service
	var serviceItem model.ServiceItem
	var req ServiceItemReq
	if err := db.Where("account_id = ?", accID).Find(&service, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	SSTime, _ := time.Parse("15:04", req.Time)
	serviceItem.SSTime = SSTime
	serviceItem.SSPrice = req.Price
	serviceItem.SSName = req.Name
	serviceItem.AccountID = accID
	if err := db.Model(&service).Association("ServiceItems").Append(&serviceItem).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"redirect": fmt.Sprintf("/admin/service/%d", service.ID),
		"data":     serviceItem,
	})
}

func ServiceItemEditViewHandler(c *Context) error {
	id := c.Param("id")
	seriveItemID := c.Param("seriveItemID")
	accID := auth.Default(c)
	db := model.DB()
	var service model.Service
	var serviceItem model.ServiceItem
	if err := db.Where("account_id = ?", accID).Find(&service, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Where("account_id = ?", accID).Find(&serviceItem, seriveItemID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.Render(http.StatusCreated, "service-item-form", echo.Map{
		"service": service,
		"detail":  serviceItem,
		"title":   "service",
		"method":  "PUT",
	})
}

func ServiceItemEditPutHandler(c *Context) error {
	id := c.Param("id")
	seriveItemID := c.Param("seriveItemID")
	accID := auth.Default(c)
	db := model.DB()
	var service model.Service
	var serviceItem model.ServiceItem
	if err := db.Where("account_id = ?", accID).Find(&service, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Where("account_id = ?", accID).Find(&serviceItem, seriveItemID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Bind(&serviceItem); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Save(&serviceItem).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"redirect": fmt.Sprintf("/admin/service/%d", service.ID),
		"data":     serviceItem,
	})
}

func ServiceItemRemoveHandler(c *Context) error {
	id := c.Param("id")
	seriveItemID := c.Param("seriveItemID")
	accID := auth.Default(c).GetAccountID()
	db := model.DB()
	var service model.Service
	var serviceItem model.ServiceItem
	if err := db.Where("account_id = ?", accID).Find(&service, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Where("account_id = ?", accID).Find(&serviceItem, seriveItemID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Delete(&serviceItem).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"data": serviceItem,
	})
}

func ServiceItemIsActiveHandler(c *Context) error {
	id := c.Param("id")
	seriveItemID := c.Param("seriveItemID")
	accID := auth.Default(c).GetAccountID()
	db := model.DB()
	var service model.Service
	var serviceItem model.ServiceItem
	isActive, err := strconv.ParseBool(c.FormValue("s_s_is_active"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Where("account_id = ?", accID).Find(&service, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Where("account_id = ?", accID).Find(&serviceItem, seriveItemID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	serviceItem.SSIsActive = isActive
	if err := db.Save(&serviceItem).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, serviceItem)
}
