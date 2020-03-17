package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IntouchOpec/base-go-echo/lib/lineapi"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"

	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

func PaymentOmiseHandler(c echo.Context) error {
	accountName := c.QueryParam("account_name")
	DocCodeTransaction := c.QueryParam("doc_code_transaction")
	LiffID := c.QueryParam("liff_id")
	var transaction model.Transaction
	var account model.Account
	db := model.DB()
	db.Where("acc_name = ?", accountName).Find(&account)
	db.Where("account_id = ? and tran_document_code = ?", account.ID, DocCodeTransaction).Find(&transaction)
	var chatChannel model.ChatChannel
	db.Find(&chatChannel, transaction.ChatChannelID)

	if transaction.TranStatus == model.TranStatusPaid {
		return c.Render(http.StatusOK, "payment-success", echo.Map{
			"accountName":        accountName,
			"DocCodeTransaction": DocCodeTransaction,
			"detail":             transaction,
			"title":              "ชำระเงินเรียบร้อยแล้ว",
		})
	}

	return c.Render(http.StatusOK, "payment-omise", echo.Map{
		"accountName":        accountName,
		"DocCodeTransaction": DocCodeTransaction,
		"detail":             transaction,
		"LiffID":             LiffID,
	})
}

const (
	OmisePublicKey = "pkey_test_5ip8fflleizk5mzvnut"
	OmiseSecretKey = "skey_test_5ip8nm6pyp7ziztxlh9"
)

type LineMassage struct {
	code        string
	account     model.Account
	transaction model.Transaction
	chatChannel model.ChatChannel
	receipt     model.Receipt
	charge      omise.Charge
}

func ChargeOmiseHandler(c echo.Context) error {
	accountName := c.QueryParam("account_name")
	DocCodeTransaction := c.QueryParam("doc_code_transaction")
	var lm LineMassage
	sqlDB := model.SqlDB()
	// get transaction
	row := sqlDB.QueryRow(`
	SELECT 
		ac.id AS account_id, acc_name, 
		cc.id AS chat_channel_id, cha_channel_secret, cha_channel_access_token, cha_address, 
		tr.id AS transaction_id, tran_document_code, tran_total, tran_line_id
	FROM transactions AS tr
	INNER JOIN chat_channels AS cc ON tr.chat_channel_id = cc.ID AND cc.deleted_at IS NULL 
	INNER JOIN accounts AS ac ON cc.account_id = ac.ID AND ac.deleted_at IS NULL AND ac.acc_name = $1
	WHERE tr.deleted_at IS NULL AND tran_document_code = $2`, accountName, DocCodeTransaction)
	err := row.Scan(
		&lm.account.ID,
		&lm.account.AccName,
		&lm.chatChannel.ID,
		&lm.chatChannel.ChaChannelSecret,
		&lm.chatChannel.ChaChannelAccessToken,
		&lm.chatChannel.ChaAddress,
		&lm.transaction.ID,
		&lm.transaction.TranDocumentCode,
		&lm.transaction.TranTotal,
		&lm.transaction.TranLineID,
	)
	if err != nil {
		fmt.Println("err", err)
		lm.code = "notFound"
		return c.JSON(http.StatusBadRequest, err)
	}

	ms, vStr, err := lm.transaction.MakeMasterBooking(sqlDB)
	if err != nil {
		lm.code = vStr
		if err := lm.sandMassage(sqlDB); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	// payment
	client, err := omise.NewClient(OmisePublicKey, OmiseSecretKey)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	token := c.FormValue("token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{})
	}
	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:   int64(lm.transaction.TranTotal * 100),
		Currency: "thb",
		Card:     token,
	}
	if err := client.Do(charge, createCharge); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	var omiseLog model.OmiseLog
	ev, err := json.Marshal(charge)
	omiseLog.Json = ev
	omiseLog.AccountID = lm.account.ID
	if err := omiseLog.Create(sqlDB); err != nil {

	}
	if charge.Status != omise.ChargeSuccessful {
		fmt.Println("err", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	lm.transaction.TranStatus = model.TranStatusPaid
	tx, err := sqlDB.Begin()
	if err := model.CreateMasterBooking(vStr, tx, ms); err != nil {
		fmt.Println(err, "Err")
	}
	payment := lm.bindPayment()
	// create
	if err := payment.Create(tx); err != nil {
		fmt.Println(err, "===1")
		tx.Rollback()
	}
	_, err = lm.transaction.UpdateStatus(tx)
	if err != nil {
		fmt.Println(err, "===2")
		tx.Rollback()
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err, "===3")
	}
	lm.code = "succes"
	if err := lm.sandMassage(sqlDB); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{})
	}

	return c.JSON(http.StatusOK, echo.Map{})
}

func (lm LineMassage) bindPayment() model.Payment {
	var payment model.Payment
	payment.PayAt = lm.charge.Created
	payment.TransactionID = lm.transaction.ID
	payment.PayStatus = model.PayStatusSuccess
	payment.PayAmount = lm.transaction.TranTotal
	payment.PayType = model.PayTypeOmise
	lm.transaction.TranStatus = model.TranStatusPaid
	return payment
}

func (lm LineMassage) sandMassage(sqlDB *sql.DB) error {
	var card string
	var text string
	bot, err := lineapi.ConnectLineBot(lm.chatChannel.ChaChannelSecret, lm.chatChannel.ChaChannelAccessToken)
	if err != nil {
		return err
	}
	switch lm.code {
	case "succes":
		reList, err := lm.transaction.GetReceipt(sqlDB)
		var list string
		if err != nil {
			return err
		}
		for _, item := range reList {
			list += fmt.Sprintf(listTemplate,
				item.Name,
				item.Price)
		}
		text = "ชำรเงินเสำเร็จ"
		card = fmt.Sprintf(receiptTemplate, lm.account.AccName, lm.chatChannel.ChaAddress, list, len(reList), lm.transaction.TranTotal, lm.transaction.TranDocumentCode)
	case "notPlace":
		card = fmt.Sprintf(cardTemplate, lm.code)
	case "notEmployee":
		card = fmt.Sprintf(cardTemplate, lm.code)
	case "notEmployeeReady":
		card = fmt.Sprintf(cardTemplate, lm.code)
	case "notPlaceReady":
		card = fmt.Sprintf(cardTemplate, lm.code)
	}
	text = lm.code

	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(card))
	if err != nil {
		fmt.Println(err, "====")
		return err
	}
	message := linebot.NewFlexMessage(text, flexContainer)
	_, err = bot.PushMessage(lm.transaction.TranLineID, message).Do()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

var cardTemplate string = `
{
	"type": "bubble",
	"hero": { "type": "image", "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/01_1_cafe.png", "size": "full", "aspectRatio": "20:13", "aspectMode": "cover" },
	"body": { "type": "box", "layout": "vertical", "contents": [
		{ "type": "text", "text": "%s", "weight": "bold", "size": "xl" }
	  ] },
	"footer": {
	  "type": "box",
	  "layout": "vertical",
	  "spacing": "sm",
	  "contents": [
		{
		  "type": "button",
		  "style": "link",
		  "height": "sm",
		  "action": {
				"type": "postback",
				"label": "booking",
				"data": "action=service"
			}
		},
		{
		  "type": "spacer",
		  "size": "sm"
		}
	  ],
	  "flex": 0
	}
  }
`

var listTemplate string = `
{ "type": "box", "layout": "horizontal",
	"contents": [
		{ "type": "text", "text": "%s", "size": "sm", "color": "#555555", "flex": 0 },
		{ "type": "text", "text": "฿ %s", "size": "sm", "color": "#111111", "align": "end"}]},`

var receiptTemplate string = `{
	"type": "bubble",
	"body": { "type": "box", "layout": "vertical", "contents": [
		{ "type": "text", "text": "RECEIPT", "weight": "bold", "color": "#1DB446", "size": "sm" },
		{ "type": "text", "text": "%s", "weight": "bold", "size": "xxl", "margin": "md" },
		{ "type": "text", "text": "%s", "size": "xs", "color": "#aaaaaa", "wrap": true },
		{
		  "type": "separator",
		  "margin": "xxl"
		},
		{ "type": "box", "layout": "vertical", "margin": "xxl", "spacing": "sm", "contents": [
			%s
			{ "type": "separator", "margin": "xxl" },
			{ "type": "box", "layout": "horizontal", "margin": "xxl", "contents": [
				{ "type": "text", "text": "ITEMS", "size": "sm", "color": "#555555" },
				{ "type": "text", "text": "%d", "size": "sm", "color": "#111111", "align": "end" }
			  ]
			},
			{ "type": "box", "layout": "horizontal", "contents": [
				{ "type": "text", "text": "TOTAL", "size": "sm", "color": "#555555" },
				{ "type": "text", "text": "฿ %f", "size": "sm", "color": "#111111", "align": "end" } ]
			}
		  ]
		},
		{ "type": "separator", "margin": "xxl" },
		{ "type": "box", "layout": "horizontal", "margin": "md",
		  "contents": [
			{ "type": "text", "text": "PAYMENT ID", "size": "xs", "color": "#aaaaaa", "flex": 0 },
			{ "type": "text", "text": "#%s", "color": "#aaaaaa", "size": "xs", "align": "end" }
		  ]
		}
	  ]
	},
	"styles": {
	  "footer": { "separator": true }
	}
  }`
