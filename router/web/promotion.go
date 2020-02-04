package web

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/jinzhu/gorm"

	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/labstack/echo"
)

// PromotionListHandler
func PromotionListHandler(c *Context) error {
	promotions := []*model.Promotion{}
	a := auth.Default(c)
	queryPar := c.QueryParams()
	page, limit := SetPagination(queryPar)
	var total int
	db := model.DB()

	filterPromotion := db.Model(&promotions).Where("account_id = ?", a.GetAccountID()).Count(&total)
	pagination := MakePagination(total, page, limit)
	filterPromotion.Limit(pagination.Record).Offset(pagination.Offset).Find(&promotions)

	err := c.Render(http.StatusOK, "promotion-list", echo.Map{
		"list":       promotions,
		"title":      "promotion",
		"pagination": pagination,
	})
	return err
}

// PromotionDetailHandler
func PromotionDetailHandler(c *Context) error {
	var chatChannels []model.ChatChannel
	promotion := model.Promotion{}
	id := c.Param("id")
	a := auth.Default(c)

	model.DB().Preload("Account").Preload("Coupons", func(db *gorm.DB) *gorm.DB {
		return db.Preload("ChatChannel")
	}).Preload("Vouchers", func(db *gorm.DB) *gorm.DB {
		return db.Preload("ChatChannel")
	}).Where("account_id = ?",
		a.GetAccountID()).Find(&promotion, id)

	if err := model.DB().Where("account_id = ?", a.GetAccountID()).Find(&chatChannels).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.Render(http.StatusOK, "promotion-detail", echo.Map{
		"detail":       promotion,
		"chatChannels": chatChannels,
		"title":        "promotion",
	})
}

type PromotionDetailForm struct {
	PromotionDetail *model.PromotionDetail
	ChatChannels    []model.ChatChannel
}
type VoucherForm struct {
	Voucher      *model.Voucher
	ChatChannels []model.ChatChannel
}

type CouponForm struct {
	Coupon       *model.Coupon
	ChatChannels []model.ChatChannel
}

func PromotionFormHandler(c *Context) error {
	promotion := model.Promotion{}
	chatChannels := []model.ChatChannel{}
	accID := auth.Default(c).GetAccountID()
	db := model.DB()
	db.Where("account_id = ?", accID).Find(&chatChannels)
	promotionTypes := []model.PromotionType{model.PromotionPromotionType, model.PromotionTypeCoupon, model.PromotionTypeVoucher}
	return c.Render(http.StatusOK, "promotion-form", echo.Map{
		"method":         "POST",
		"chatChannels":   chatChannels,
		"detail":         promotion,
		"title":          "promotion",
		"promotionTypes": promotionTypes,
		"PromotionForm":  &PromotionDetailForm{PromotionDetail: &model.PromotionDetail{}, ChatChannels: chatChannels},
		"VoucherForm":    &VoucherForm{Voucher: &model.Voucher{}, ChatChannels: chatChannels},
		"CouponForm":     &CouponForm{Coupon: &model.Coupon{}, ChatChannels: chatChannels},
	})
}

type PromotionForm struct {
	Title         string `form:"title"`
	PromotionType string `form:"promotion_type"`
	Discount      int    `form:"discount"`
	Amount        int    `form:"amount"`
	Code          string `form:"code"`
	Name          string `form:"name"`
}

func PromotionPostHandler(c *Context) error {
	file := c.FormValue("file")
	accID := auth.Default(c).GetAccountID()
	promotion := PromotionForm{}
	if err := c.Bind(&promotion); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	ctx := context.Background()
	imagePath, err := lib.UploadGoolgeStorage(ctx, file, "images/promotion/")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	promotionModel := model.Promotion{
		PromTitle:    promotion.Title,
		PromType:     promotion.PromotionType,
		PromDiscount: promotion.Discount,
		PromAmount:   promotion.Amount,
		PromCode:     promotion.Code,
		PromName:     promotion.Name,
		PromImage:    imagePath,
		AccountID:    accID,
	}
	if err := model.DB().Create(&promotionModel).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// promotionModel.SavePromotion()
	return c.JSON(http.StatusCreated, echo.Map{
		"data":     promotionModel,
		"redirect": fmt.Sprintf("/admin/promotion/%d", promotionModel.ID),
	})
}

type reqSubPromotion struct {
	StartDate string `form:"start_date" json:"start_date"`
	EndDate   string `form:"end_date" json:"end_date"`
	Amount    string `form:"amount" json:"amount"`
	Condition string `form:"condition" json:"condition"`
}

type Timestamp time.Time

func (t *Timestamp) UnmarshalParam(src string) error {
	ts, err := time.Parse(time.RFC3339, src)
	*t = Timestamp(ts)
	return err
}

func PromotionCreateDetailHandler(c *Context) error {
	id := c.Param("id")
	var promotion model.Promotion
	var chatChannel model.ChatChannel
	var req reqSubPromotion
	accID := auth.Default(c).GetAccountID()
	db := model.DB()
	if err := db.Find(&promotion, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	chatChannelID, _ := strconv.ParseUint(c.FormValue("chat_channel_id"), 10, 32)
	startDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	amount, err := strconv.Atoi(req.Amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	switch promotion.PromType {
	case "Coupon":
		db.Model(&promotion).Association("Coupons").Append(&model.Coupon{
			PromStartDate: startDate,
			PromEndDate:   endDate,
			PromAmount:    amount,
			PromCondition: req.Condition,
			ChatChannelID: uint(chatChannelID),
			AccountID:     accID,
		})
	case "Voucher":
		db.Model(&promotion).Association("Vouchers").Append(&model.Voucher{
			PromStartDate: startDate,
			PromEndDate:   endDate,
			PromAmount:    amount,
			PromCondition: req.Condition,
			ChatChannelID: uint(chatChannelID),
			AccountID:     accID,
		})
	}

	if err := db.Find(&chatChannel, chatChannelID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"data":         promotion,
		"chat_channel": chatChannel,
	})
}

func PromotionEditHandler(c *Context) error {
	promotion := model.Promotion{}
	id := c.Param("id")
	a := auth.Default(c)

	model.DB().Preload("ChatChannels").Where("account_id = ?", a.User.GetAccountID()).Find(&promotion, id)
	return c.Render(http.StatusOK, "promotion-form", echo.Map{
		"method": "PUT",
		"detail": promotion,
		"title":  "promotion",
	})
}

func PromotionEditPutHandler(c *Context) error {
	id := c.Param("id")
	a := auth.Default(c)
	var err error
	imagePath := c.FormValue("image")
	if imagePath == "" {
		file := c.FormValue("file")
		ctx := context.Background()
		imagePath, err = lib.UploadGoolgeStorage(ctx, file, "images/promotion/")
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	promotion := PromotionForm{}
	promotionModel := model.Promotion{}
	if err := c.Bind(&promotion); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	db := model.DB()
	if err := db.Where("account_id = ?", a.User.GetAccountID()).Find(&promotionModel, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	promotionModel.PromTitle = promotion.Title
	promotionModel.PromType = promotion.PromotionType
	promotionModel.PromDiscount = promotion.Discount
	// promotionModel.PromAmount = promotion.Amount
	// promotionModel.PromCode = promotion.Code
	// promotionModel.PromName = promotion.Name
	// promotionModel.PromStartDate = promotion.StartDate
	// promotionModel.PromEndDate = promotion.EndDate
	// promotionModel.PromCondition = promotion.Condition
	promotionModel.PromImage = imagePath
	if err := db.Save(&promotionModel).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"detail":   promotionModel,
		"redirect": fmt.Sprintf("/admin/promotion/%s", id),
	})
}

func PromotionChannelFormHandler(c *Context) error {
	chatChannels := []*model.ChatChannel{}
	a := auth.Default(c)

	if err := model.DB().Where("account_id = ?", a.User.GetAccountID()).Find(&chatChannels).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.Render(http.StatusOK, "promotion-chat-channel-form", echo.Map{
		"method":       "PUT",
		"chatChannels": chatChannels,
		"title":        "promotion",
		"mode":         "Create",
	})
}

func PromotionChannelAddHandler(c *Context) error {
	a := auth.Default(c)
	id := c.Param("id")
	pro := model.Promotion{}
	chatChannel := model.ChatChannel{}
	chatChannelID := c.FormValue("chat_channel_id")
	db := model.DB()
	if err := db.Where("account_id = ?", a.User.GetAccountID()).Find(&chatChannel, chatChannelID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := db.Where("account_id = ?", a.User.GetAccountID()).Find(&pro, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Model(&pro).Association("ChatChannels").Append(&chatChannel).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	redirect := fmt.Sprintf("/admin/promotion/%d", pro.ID)
	return c.JSON(http.StatusCreated, echo.Map{
		"redirect": redirect,
	})
}

func PromotionRemoveHandler(c *Context) error {
	id := c.Param("id")

	promotion, err := model.DeletePromotion(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, promotion)
}

func PromotionAddRegisterlHandler(c *Context) error {
	id := c.Param("id")
	voucher := model.Voucher{}
	chatChannel := model.ChatChannel{}
	db := model.DB()
	chatChannelID := c.FormValue("chat_channel_id")

	if err := db.Find(&chatChannel, chatChannelID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := db.Find(&voucher, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	chatChannel.VoucherID = voucher.ID
	if err := db.Save(&chatChannel).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, chatChannel)
}

func PromotionDeleteImageHandler(c *Context) error {
	id := c.Param("id")
	promotion := model.Promotion{}
	accID := auth.Default(c).GetAccountID()
	db := model.DB()
	if err := db.Where("account_id = ? ", accID).Find(&promotion, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	ctx := context.Background()
	if _, err := lib.RemoveFileGoolgeStorage(ctx, "triple-t", promotion.PromImage); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	promotion.PromImage = ""
	if err := db.Save(&promotion).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"detail": promotion,
	})
}
