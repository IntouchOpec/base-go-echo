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

	"github.com/jinzhu/gorm"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	guuid "github.com/google/uuid"
	"github.com/labstack/echo"
)

// ProductListHandler
func ProductListHandler(c *Context) error {
	products := []*model.Product{}
	a := auth.Default(c)
	model.DB().Preload("ChatChannel", func(db *gorm.DB) *gorm.DB {
		return db.Where("chat_channel_id = ?", a.User.GetAccountID())
	}).Preload("SubProduct").Find(&products)
	err := c.Render(http.StatusOK, "product-list", echo.Map{
		"list":  products,
		"title": "product",
	})
	return err
}

// ProductDetailHandler
func ProductDetailHandler(c *Context) error {
	product := model.Product{}
	id := c.Param("id")
	a := auth.Default(c)
	model.DB().Preload("Account").Preload("SubProducts").Preload("ChatChannels").Where("account_id = ? ", a.User.GetAccountID()).Find(&product, id)
	err := c.Render(http.StatusOK, "product-detail", echo.Map{
		"detail": product,
		"title":  "product",
	})
	return err
}

func ProductCreateHandler(c *Context) error {
	product := model.Product{}
	csrfValue := c.Get("_csrf")

	err := c.Render(http.StatusOK, "product-form", echo.Map{
		"detail": product,
		"title":  "product",
		"_csrf":  csrfValue,
	})
	return err
}

func ProductEditHandler(c *Context) error {
	product := model.Product{}
	id := c.Param("id")
	a := auth.Default(c)
	model.DB().Preload("Account").Preload("SubProducts").Preload("ChatChannels").Where("account_id = ? ", a.User.GetAccountID()).Find(&product, id)
	err := c.Render(http.StatusOK, "product-form", echo.Map{
		"detail": product,
		"title":  "product",
	})
	return err
}

func ProductDeleteHandler(c *Context) error {
	id := c.Param("id")
	pro := model.DeleteProductByID(id)
	err := c.JSON(http.StatusOK, pro)
	return err
}

func SubProductCreateHandler(c *Context) error {
	messageTypes := []linebot.MessageType{linebot.MessageTypeText, linebot.MessageTypeImage, linebot.MessageTypeVideo, linebot.MessageTypeAudio, linebot.MessageTypeFile, linebot.MessageTypeLocation, linebot.MessageTypeSticker, linebot.MessageTypeTemplate, linebot.MessageTypeImagemap, linebot.MessageTypeFlex}

	sunProduct := model.SubProduct{}
	err := c.Render(http.StatusOK, "sub-product-form", echo.Map{
		"detail":       sunProduct,
		"title":        "product",
		"messageTypes": messageTypes,
	})
	return err
}

type ProductForm struct {
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

func ProductPostHandler(c *Context) error {
	product := ProductForm{}
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
	fileName := fileNameBase + "." + fm
	err = ioutil.WriteFile(fileName, buff.Bytes(), 0644)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	a := auth.Default(c)
	productModel := model.Product{
		Name:      product.Name,
		Detail:    product.Detail,
		Price:     product.Price,
		Image:     fmt.Sprintf("%s.%s", u, fm),
		AccountID: a.User.GetAccountID(),
	}
	productModel.SaveProduct()

	return c.JSON(http.StatusCreated, productModel)
}

type SubProductForm struct {
	Start  string `form:"start" json:"start"`
	End    string `form:"end" json:"end"`
	Day    int    `form:"day" json:"day"`
	Amount int    `form:"amount" json:"amount"`
}

func SubProductPostHandler(c *Context) error {
	id := c.Param("id")
	product := model.Product{}
	db := model.DB()
	subProductFrom := SubProductForm{}
	if err := c.Bind(&subProductFrom); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	db.Find(&product, id)
	db.Model(&product).Association("SubProducts").Append(&model.SubProduct{Start: subProductFrom.Start, End: subProductFrom.End, Day: subProductFrom.Day, Amount: subProductFrom.Amount})
	return c.JSON(http.StatusCreated, product)
}

func SubProductEditHandler(c *Context) error {
	sunProduct := model.SubProduct{}
	id := c.Param("id")
	a := auth.Default(c)
	messageTypes := []linebot.MessageType{linebot.MessageTypeText, linebot.MessageTypeImage, linebot.MessageTypeVideo, linebot.MessageTypeAudio, linebot.MessageTypeFile, linebot.MessageTypeLocation, linebot.MessageTypeSticker, linebot.MessageTypeTemplate, linebot.MessageTypeImagemap, linebot.MessageTypeFlex}
	model.DB().Preload("Product", func(db *gorm.DB) *gorm.DB {
		return db.Preload("ChatChannels").Where("account_id = ? ", a.User.GetAccountID())
	}).Find(&sunProduct, id)
	err := c.Render(http.StatusOK, "sub-product-form", echo.Map{
		"detail":       sunProduct,
		"title":        "product",
		"messageTypes": messageTypes,
	})
	return err
}

func SubProductDeleteHandler(c *Context) error {
	id := c.Param("id")
	subProduct := model.DeleteSubProduct(id)
	return c.JSON(http.StatusOK, subProduct)
}

func ProductChatChannelViewHandler(c *Context) error {
	chatChannels := []*model.ChatChannel{}
	a := auth.Default(c)
	model.DB().Preload("Account").Where("account_id = ?", a.User.GetAccountID()).Find(&chatChannels)

	err := c.Render(http.StatusOK, "product-chat-channel-form", echo.Map{
		"list_chat_channel": chatChannels,
		"title":             "product",
	})
	return err
}

func ProductChatChannelPostHandler(c *Context) error {
	id := c.QueryParam("id")
	chatChannelID := c.FormValue("chat_channel_id")
	product := model.Product{}
	chatChannel := model.ChatChannel{}
	db := model.DB()

	if err := db.Find(&chatChannel, chatChannelID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := db.Find(&product, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := db.Model(&product).Association("ChatChannels").Append(&chatChannel).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, product)
}
