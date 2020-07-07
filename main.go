package main

import (
	"order/application"
	"order/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	app.GET("/users/:id/orders", handlers.FindOrdersForUser)
	err := app.Run(":8000")
	if err != nil {
		panic(err)
	}
	defer application.CloseDB()
}
