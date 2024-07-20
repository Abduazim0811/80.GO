package api

import (
	"log"

	"80.GO/api/handler"
	"80.GO/internal/mongodb"
	"80.GO/internal/rabbitmq"
	"github.com/gin-gonic/gin"
)

func Routes() {
	router := gin.Default()
	db, err := mongodb.NewOrder()
	if err != nil {
		log.Fatal(err)
	}
	orderhandler := handler.NewOrderHandler(db)
	router.POST("/orders", orderhandler.CreateOrder)
	router.GET("/orders/:id", orderhandler.GetTasksbyID)

	go rabbitmq.ConsumeOrders()

	router.Run(":7777")
}
