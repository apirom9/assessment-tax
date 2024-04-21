package main

import (
	"fmt"

	"github.com/apirom9/assessment-tax/postgres"
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

	store, err := postgres.NewPostgres()
	if err != nil {
		fmt.Printf("Unable to create store DB, error: %v", err)
		return
	}

	handler := tax.Handler{Store: store}

	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/tax/calculations", handler.CalculateTax)
	e.POST("/tax/calculations/upload-csv", handler.CalculateTaxCsv)
	e.POST("/admin/deductions/personal", handler.UpdatePersonalDeduction)
	e.Logger.Fatal(e.Start(":1323"))
}
