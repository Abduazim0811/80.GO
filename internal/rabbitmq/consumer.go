package rabbitmq

import (
	"encoding/json"
	"log"

	"80.GO/internal/models"
	"80.GO/internal/mongodb"
	"github.com/streadway/amqp"
)

func ConsumeOrders() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    
		false, 
		false, 
		true,  
		false, 
		nil,   
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(
		q.Name,   
		"",       
		"orders", 
		false,
		nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name, 
		"",     
		true,   
		false,  
		false,  
		false,  
		nil,    
	)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var order models.Order
			err := json.Unmarshal(d.Body, &order)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
				continue
			}

			// Mongodb ga saqlash
			mongodb, err := mongodb.NewOrder()
			if err != nil {
				log.Fatal(err)
			}
			if err = mongodb.CreateOrderMongoDb(order); err != nil{
				log.Printf("Error saving task to MongoDB: %s", err)
                continue
			}

			log.Printf("Buyurtma olindi va qayta ishlanmoqda: %+v", order)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
