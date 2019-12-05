package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IntouchOpec/base-go-echo/lib"
	"github.com/IntouchOpec/base-go-echo/module/auth"

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
	promotion := model.Promotion{}
	id := c.Param("id")
	a := auth.Default(c)

	err := model.DB().Preload("Account").Preload("Customers").Preload("ChatChannels").Where("account_id = ?",
		a.User.GetAccountID()).Find(&promotion, id)
	if err != nil {
		fmt.Println(err, "===")
	}
	sumCustomer := len(promotion.Customers)
	return c.Render(http.StatusOK, "promotion-detail", echo.Map{
		"detail":      promotion,
		"title":       "promotion",
		"sumCustomer": sumCustomer,
	})
}

func PromotionFormHandler(c *Context) error {
	promotion := model.Promotion{}
	promotionTypes := []model.PromotionType{model.PromotionPromotionType, model.PromotionTypeCoupon, model.PromotionTypeVoucher}
	return c.Render(http.StatusOK, "promotion-form", echo.Map{"method": "PUT",
		"detail":         promotion,
		"title":          "promotion",
		"promotionTypes": promotionTypes,
	})
}

type PromotionForm struct {
	Title         string    `form:"title"`
	PromotionType string    `form:"promotion_type"`
	Discount      int       `form:"discount"`
	Amount        int       `form:"amount"`
	Code          string    `form:"code"`
	Name          string    `form:"name"`
	StartDate     time.Time `form:"start_time"`
	EndDate       time.Time `form:"end_time"`
	Condition     string    `form:"condition"`
}

func PromotionPostHandler(c *Context) error {
	file := c.FormValue("file")
	a := auth.Default(c)
	promotion := PromotionForm{}
	if err := c.Bind(&promotion); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	fileUrl, _, err := lib.UploadteImage(file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	promotionModel := model.Promotion{
		PromTitle:     promotion.Title,
		PromType:      promotion.PromotionType,
		PromDiscount:  promotion.Discount,
		PromAmount:    promotion.Amount,
		PromCode:      promotion.Code,
		PromName:      promotion.Name,
		PromStartDate: promotion.StartDate,
		PromEndDate:   promotion.EndDate,
		PromCondition: promotion.Condition,
		PromImage:     fileUrl,
		AccountID:     a.User.GetAccountID(),
	}

	promotionModel.SavePromotion()
	return c.JSON(http.StatusCreated, echo.Map{
		"data":     promotionModel,
		"redirect": fmt.Sprintf("/admin/promotion/%d", promotionModel.ID),
	})
}

func PromotionEditHandler(c *Context) error {
	promotion := model.Promotion{}
	id := c.Param("id")
	a := auth.Default(c)

	model.DB().Preload("Account").Preload("services").Preload("Customers").Preload("ChatChannels").Where("account_id = ?", a.User.GetAccountID()).Find(&promotion, id)
	return c.Render(http.StatusOK, "promotion-form", echo.Map{"method": "PUT",
		"detail": promotion,
		"title":  "promotion",
	})
}

func PromotionChannelFormHandler(c *Context) error {
	chatChannels := []*model.ChatChannel{}
	a := auth.Default(c)

	if err := model.DB().Where("account_id = ?", a.User.GetAccountID()).Find(&chatChannels).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.Render(http.StatusOK, "promotion-chat-channel-form", echo.Map{"method": "PUT",
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
	fmt.Println(id, "id")
	if err := db.Where("account_id = ?", a.User.GetAccountID()).Find(&chatChannel, chatChannelID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := db.Where("account_id = ?", a.User.GetAccountID()).Find(&pro, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := db.Model(&pro).Association("ChatChannels").Append(&chatChannel).Error; err != nil {
		fmt.Println(chatChannel)
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
	promotion := model.Promotion{}
	chatChannel := model.ChatChannel{}
	db := model.DB()
	chatChannelID := c.FormValue("chat_channel_id")

	if err := db.Find(&chatChannel, chatChannelID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := db.Find(&promotion, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	chatChannel.PromotionID = promotion.ID
	if err := db.Save(&chatChannel).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, chatChannel)
}
