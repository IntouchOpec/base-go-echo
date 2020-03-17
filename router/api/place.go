package api

// func GetPlaceList(c echo.Context) error {
// 	// a :=
// 	places, err := model.GetPlaceList(1)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, err)
// 	}
// 	return c.JSON(http.StatusOK, places)
// }

// func GetPlaceDetail(c echo.Context) error {
// 	id := c.Param("id")
// 	place, err := model.GetPlaceDetail(id, 1)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, err)
// 	}
// 	return c.JSON(http.StatusOK, place)
// }

// func CreatePlace(c echo.Context) error {
// 	pla := model.Place{}
// 	if err := c.Bind(&pla); err != nil {
// 		return c.JSON(http.StatusBadRequest, err)
// 	}
// 	if err := pla.CreatePlace(); err != nil {
// 		return c.JSON(http.StatusBadRequest, err)
// 	}
// 	return c.JSON(http.StatusCreated, pla)
// }

// func UpdatePlace(c echo.Context) error {
// 	id := c.Param("id")
// 	pla, err := model.GetPlaceDetail(id, 1)
// 	if err := c.Bind(&pla); err != nil {
// 		return c.NoContent(http.StatusBadRequest)
// 	}

// 	err = pla.Update()
// 	if err != nil {
// 		return c.NoContent(http.StatusInternalServerError)
// 	}
// 	return c.JSON(http.StatusOK, pla)
// }

// func DeletePlace(c echo.Context) error {
// 	id := c.Param("id")
// 	pla, err := model.DeletePlaceByID(id)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, err)
// 	}

// 	return c.JSON(http.StatusOK, pla)
// }
