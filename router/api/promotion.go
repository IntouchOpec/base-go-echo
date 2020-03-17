package api

// func CreatePromotion(c echo.Context) error {

// 	promotion := model.Promotion{}

// 	if err := c.Bind(&promotion); err != nil {
// 		return c.JSON(http.StatusBadRequest, err)
// 	}
// 	promotion.SavePromotion()
// 	return c.JSON(http.StatusOK, promotion)
// }

// func GetPromotionList(c echo.Context) error {
// 	chatChannelID := c.Param("chatChannelID")
// 	chatChannelIDInt, _ := strconv.Atoi(chatChannelID)
// 	promtions := model.GetPromotionList(chatChannelIDInt)
// 	return c.JSON(http.StatusOK, promtions)
// }

// func GetPromotion(c echo.Context) error {
// 	id := c.Param("id")
// 	idInt, _ := strconv.Atoi(id)
// 	promotion := model.GetPromotion(idInt)
// 	return c.JSON(http.StatusOK, promotion)
// }

// func UpdatePromotion(c echo.Context) error {
// 	id := c.Param("id")
// 	idInt, _ := strconv.Atoi(id)

// 	promotion := model.Promotion{}
// 	if err := c.Bind(&promotion); err != nil {
// 		return c.NoContent(http.StatusBadRequest)
// 	}
// 	promotion.UpdatePromotion(idInt)
// 	return c.JSON(http.StatusOK, promotion)
// }

// func DeletePromotion(c echo.Context) error {
// 	id := c.Param("id")
// 	promotion, err := model.DeletePromotion(id)
// 	if err != nil {
// 		return c.JSON(http.StatusOK, err)
// 	}
// 	return c.JSON(http.StatusOK, promotion)
// }
