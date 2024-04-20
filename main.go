package main

import (
	"github.com/apirom9/assessment-tax/tax"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/apirom9/assessment-tax/docs"
)

// @title			Tax API
// @version		1.0
// @description	Tax API
// @host			localhost:1323
func main() {
	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/tax/calculations", tax.CalculateTax)
	e.POST("/tax/calculations/upload-csv", tax.CalculateTaxCsv)
	e.Logger.Fatal(e.Start(":1323"))
}
