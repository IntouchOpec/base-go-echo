package web

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

// LIFFloginHandler
func LIFFRegisterHandler(c echo.Context) error {
	lineID := c.Param("lineID")
	chatChannel := model.ChatChannel{}
	customerTypes := []model.CustomerType{}
	db := model.DB()
	if err := db.Where("cha_line_id = ?", lineID).Find(&chatChannel).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	custo := model.Customer{}
	if err := db.FirstOrCreate(&custo).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	fmt.Println(chatChannel.AccountID)
	db.Where("account_id = ?", chatChannel.AccountID).Find(&customerTypes)
	APIRegister := fmt.Sprintf("https://web.%s/register/%s", Conf.Server.Domain, lineID)
	// APIRegister := fmt.Sprintf("https://%s/register/%s", "586f1140.ngrok.io", lineID)

	return c.Render(http.StatusOK, "register", echo.Map{
		"web":           APIRegister,
		"customerTypes": customerTypes,
	})
}

type LineReqRegister struct {
	UserID      string `json:"userId"`
	FullName    string `json:"fullName"`
	DisplayName string `json:"displayName"`
	PictureURL  string `json:"pictureUrl"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Type        string `json:"type"`
	AccessToken string `json:"accessToken"`
}

func LIIFRegisterSaveCustomer(c echo.Context) error {
	lineID := c.Param("lineID")

	chatChannel := model.ChatChannel{}
	req := LineReqRegister{}
	voucherCustomer := model.VoucherCustomer{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	db := model.DB()

	if err := db.Preload("Voucher", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Promotion")
	}).Where("cha_line_ID = ?", lineID).Find(&chatChannel).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	custo := model.Customer{CusLineID: req.UserID, AccountID: chatChannel.AccountID}
	bot, err := lib.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	_, err = custo.UpdateCustomerByAtt(req.PictureURL, req.DisplayName, req.Email, req.Phone, req.FullName, req.Type)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Preload("Voucher", "promotion_id = ?",
		chatChannel.Voucher.PromotionID).Where("customer_id = ?", custo.ID).Find(&voucherCustomer).Error; err != nil {
		voucherCustomer.VoucherID = chatChannel.VoucherID
		voucherCustomer.CustomerID = custo.ID
		voucherCustomer.AccountID = chatChannel.AccountID
		db.Create(&voucherCustomer)
	}

	if voucherCustomer.Voucher == nil {
		voucherCustomer.VoucherID = chatChannel.VoucherID
		voucherCustomer.CustomerID = custo.ID
		voucherCustomer.AccountID = chatChannel.AccountID
		db.Create(&voucherCustomer)
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(voucherTemplate(chatChannel.Voucher.Promotion, chatChannel.Voucher)))
	flexMessage := linebot.NewFlexMessage(chatChannel.Voucher.Promotion.PromName, flexContainer)
	// flexMessage := linebot.NewTextMessage(chatChannel.ChaWelcomeMessage)
	if _, err = bot.ReplyMessage(req.AccessToken, flexMessage).Do(); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, custo)
}

func voucherTemplate(promotion *model.Promotion, voucher *model.Voucher) string {
	StartDateStr := voucher.PromStartDate.Format("02-12-2006")
	EndDateStr := voucher.PromEndDate.Format("02-12-2006")
	temp := fmt.Sprintf(`{
		"type": "bubble",
		"hero": { "type": "image", "url": "%s", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover" },
		"body": { "type": "box", "layout": "vertical", "contents": [
			{ "type": "text", "text": "%s", "weight": "bold", "size": "xl"},
			{ "type": "box", "layout": "vertical", "margin": "lg", "spacing": "sm", "contents": [
				{ "type": "box", "layout": "baseline", "contents": [
					{ "type": "text", "text": "Date", "color": "#aaaaaa", "size": "sm", "flex": 1 },
					{ "type": "text", "text": "%s", "wrap": true, "color": "#666666", "size": "sm", "flex": 5 } ]
				},
				{ "type": "box", "layout": "baseline", "spacing": "sm", "contents": [
					{ "type": "text", "text": "exp", "color": "#aaaaaa", "flex": 1 },
					{ "type": "text", "text": "%s", "wrap": true, "color": "#666666", "size": "sm", "flex": 5 } ]
				},
				{"type": "text", "margin": "lg", "text": "%s", "align": "center"},
				{"type": "button", "style": "secondary", "action": { "type": "uri", "label": "%s", "uri": "https://web.linecorp.com" }
				}
			  ]
			}
		  ]
		},
		"footer": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
			{ "type": "button", "style": "link", "height": "sm",
			  "action": { "type": "uri", "label": "เงื่อนไขการใช้", "uri": "https://web.linecorp.com" } }
		  ],
		  "flex": 0
		}
	  }`, "https://web."+Conf.Server.DomainWeb+promotion.PromImage, promotion.PromTitle, StartDateStr, EndDateStr, voucher.PromCondition, promotion.PromCode)
	return temp
}
