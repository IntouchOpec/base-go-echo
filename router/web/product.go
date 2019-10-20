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
	uuid "github.com/satori/go.uuid"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
)

// ProductListHandler
func ProductListHandler(c *Context) error {
	products := []*model.Product{}
	a := auth.Default(c)
	model.DB().Preload("ChatChannel", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Account", "name = ?", a.User.GetAccountID())
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
	u1,  := uuid.Must(uuid.NewV4())
	fileNameBase := "public/assets/images/%s"
	fileNameBase = fmt.Sprintf(fileNameBase, u1)
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
		Image:     fmt.Sprintf("%s.%s", u1, fm),
		AccountID: a.User.GetAccountID(),
	}
	productModel.SaveProduct()

	return c.JSON(http.StatusCreated, productModel)
}

func SubProductPostHandler(c *Context) error {
	// file, err := c.FormFile("file")
	// if err != nil {
	// 	return err
	// }
	// src, err := file.Open()
	// if err != nil {
	// 	return err
	// }
	// defer src.Close()

	// dst, err := os.Create(file.Filename)
	// if err != nil {
	// 	return err
	// }
	// defer dst.Close()

	// if _, err = io.Copy(dst, src); err != nil {
	// 	return err
	// }
	// fmt.Println(file.Filename, "====")

	return c.JSON(http.StatusCreated, "")
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
	return c.JSON(http.StatusCreated, "")
}
