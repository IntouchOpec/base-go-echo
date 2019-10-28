package web

import (
	"fmt"
	"log"
	"net/http"

	// . "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

// LIFFloginHandler
func LIFFRegisterHandler(c echo.Context) error {
	lineID := c.Param("lineID")
	chatChannel := model.ChatChannel{}
	if err := model.DB().Where("line_id = ?", lineID).Find(&chatChannel).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	// custo := model.Customer{}
	// if err := model.DB().FirstOrCreate(&custo).Error; err != nil {
	// 	return err
	// }
	// APIRegister := fmt.Sprintf("https://%s/register/%s", Conf.Server.DomainLineChannel, lineID)
	err := c.Render(http.StatusOK, "register", echo.Map{
		// "web": APIRegister,
	})
	return err
}

type LineReqRegister struct {
	UserID      string `json:"userId"`
	FullName    string `json:"fullName"`
	DisplayName string `json:"displayName"`
	PictureURL  string `json:"pictureUrl"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	AccessToken string `json:"accessToken"`
}

func LIIFRegisterSaveCustomer(c echo.Context) error {
	lineID := c.Param("lineID")

	chatChannel := model.ChatChannel{}
	req := LineReqRegister{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := model.DB().Where("line_ID = ?", lineID).Find(&chatChannel).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	custo := model.Customer{LineID: req.UserID, AccountID: chatChannel.AccountID}
	// pictureURL string, displayName string, email string, phoneNumber string
	bot, err := lib.ConnectLineBot(chatChannel.ChannelSecret, chatChannel.ChannelAccessToken)

	custo.UpdateCustomerByAtt(req.PictureURL, req.DisplayName, req.Email, req.Phone, req.FullName)
	check := ValidateVoucher(custo.Promotions)
	if check {
		return c.NoContent(http.StatusBadRequest)
	}
	promotion := model.Promotion{}

	if err := model.DB().Where("name = ?", "register_voucher").Find(&promotion).Error; err != nil {
		return err
	}

	if err := model.DB().Model(&custo).Association("Promotions").Append(promotion).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	temp := VoucherTemplate(&promotion)
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(temp))

	if err != nil {
		log.Print(err)
		return err
	}

	flexMessage := linebot.NewFlexMessage(chatChannel.WelcomeMessage, flexContainer)

	if _, err = bot.PushMessage(req.UserID, flexMessage).Do(); err != nil {
		log.Print(err)
		return err
	}

	return c.JSON(http.StatusOK, custo)
}

func ValidateVoucher(promotions []*model.Promotion) bool {
	for index := 0; index < len(promotions); index++ {
		if promotions[index].Name == "register_voucher" {
			return true
		}
	}
	return false
}

func VoucherTemplate(promotion *model.Promotion) string {
	settings := promotion.GetSettingPromotion([]string{"displayText", "TitleName"})
	StartDateStr := promotion.StartDate.Format("02-12-2006")
	EndDateStr := promotion.EndDate.Format("02-12-2006")
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
				{"type": "button", "style": "secondary", "action": { "type": "uri", "label": "%s", "uri": "https://linecorp.com" }
				}
			  ]
			}
		  ]
		},
		"footer": { "type": "box", "layout": "vertical", "spacing": "sm", "contents": [
			{ "type": "button", "style": "link", "height": "sm",
			  "action": { "type": "uri", "label": "เงื่อนไขการใช้", "uri": "https://linecorp.com" } }
		  ],
		  "flex": 0
		}
	  }`, promotion.Image, settings["TitleName"], StartDateStr, EndDateStr, settings["displayText"], promotion.Code)
	return temp
}
