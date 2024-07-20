package rabbitmq

import (
	"encoding/json"

	"80.GO/internal/models"
	"github.com/streadway/amqp"
)

func PublishOrder(order models.Order) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"orders", 
		"fanout", 
		true,     
		false,    
		false,    
		false,    
		nil,      
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"orders", 
		"",       
		false,    
		false,    
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}
