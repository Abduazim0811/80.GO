package handler

import (
	"net/http"

	"80.GO/internal/models"
	"80.GO/internal/mongodb"
	"80.GO/internal/rabbitmq"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderHandler struct {
	db *mongodb.OrderMongoDb
}

func NewOrderHandler(db *mongodb.OrderMongoDb) *OrderHandler {
	return &OrderHandler{db: db}
}

func (o *OrderHandler) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := rabbitmq.PublishOrder(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (o *OrderHandler) GetTasksbyID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	orders, err := o.db.GetOrderMongoDb(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, orders)
}
