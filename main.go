package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/apirom9/assessment-tax/postgres"
	"github.com/apirom9/assessment-tax/tax"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	docs "github.com/apirom9/assessment-tax/docs"
)

// @title			Tax API
// @version		1.0
// @description	Tax API
func main() {

	registerGracefulShutdown()

	dbUrl := os.Getenv("DATABASE_URL")
	store, err := postgres.NewPostgres(dbUrl)
	if err != nil {
		fmt.Printf("Unable to create store DB, error: %v", err)
		return
	}

	handler := tax.Handler{Store: store}

	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/tax/calculations", handler.CalculateTax)
	e.POST("/tax/calculations/upload-csv", handler.CalculateTaxCsv)

	adminUserName := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	g := e.Group("")
	g.Use(middleware.BasicAuth(func(user, password string, ctx echo.Context) (bool, error) {
		if user == adminUserName && password == adminPassword {
			return true, nil
		}
		return false, nil
	}))
	g.POST("/admin/deductions/personal", handler.UpdatePersonalDeduction)
	g.POST("/admin/deductions/k-receipt", handler.UpdateKReceipt)

	port := os.Getenv("PORT")
	docs.SwaggerInfo.Host = "localhost:" + port
	e.Logger.Fatal(e.Start(":" + port))
}

func registerGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("shutting down the server")
		os.Exit(0)
	}()
}
