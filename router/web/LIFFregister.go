package web

import (
	"fmt"
	"net/http"
	"strconv"

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
	var custo model.Customer
	var chatChannel model.ChatChannel
	var voucherCustomer model.VoucherCustomer
	var req LineReqRegister

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	db := model.DB()

	if err := db.Preload("Voucher", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Promotion")
	}).Where("cha_line_ID = ?", lineID).Find(&chatChannel).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	bot, err := lib.ConnectLineBot(chatChannel.ChaChannelSecret, chatChannel.ChaChannelAccessToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := db.Where("cus_line_id = ? and account_id = ?", req.UserID, chatChannel.AccountID).Find(&custo).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	custo.CusPictureURL = req.PictureURL
	custo.CusDisplayName = req.DisplayName
	custo.CusEmail = req.Email
	custo.CusFullName = req.FullName
	custo.CusPhoneNumber = req.Phone
	u64, _ := strconv.ParseUint(req.Type, 10, 32)
	custo.CustomerTypeID = uint(u64)
	if err := db.Save(&custo).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if chatChannel.Voucher != nil {
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
		if _, err = bot.PushMessage(req.UserID, flexMessage).Do(); err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusOK, custo)
	}
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(welcomeTemplate()))
	flexMessage := linebot.NewFlexMessage("welcome", flexContainer)
	if _, err = bot.PushMessage(req.UserID, flexMessage).Do(); err != nil {
		fmt.Println(err)

		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, custo)
}

func welcomeTemplate() string {
	return `{
		"type": "bubble",
		"hero": {
		  "type": "image",
		  "url": "https://media3.mensxp.com/media/content/2019/Aug/here-are-the-top-gadgets-you-can-gift-your-sister-this-raksha-bandhan-1400x653-1565263792_1100x513.jpg",
		  "size": "full",
		  "aspectRatio": "20:13",
		  "aspectMode": "cover",
		  "action": {
			"type": "uri",
			"uri": "http://linecorp.com/"
		  }
		},
		"body": {
		  "type": "box",
		  "layout": "vertical",
		  "contents": [
			{
			  "type": "text",
			  "text": "ขอบคุณ ที่มาเป็นเพื่อนกัน",
			  "weight": "bold",
			  "size": "xl"
			}
		  ]
		}
	  }`
}

func voucherTemplate(promotion *model.Promotion, voucher *model.Voucher) string {
	StartDateStr := voucher.PromStartDate.Format("2006-01-02")
	EndDateStr := voucher.PromEndDate.Format("2006-01-02")
	if voucher.PromCondition == "" {
		voucher.PromCondition = "-"
	}
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
	  }`, fmt.Sprintf("https://web.%s/files?path=%s", Conf.Server.Domain, promotion.PromImage), promotion.PromTitle, StartDateStr, EndDateStr, voucher.PromCondition, promotion.PromCode)
	return temp
}
