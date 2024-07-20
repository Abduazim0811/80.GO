package models

type Order struct{
	ID string 			`json:"id" bson:"id"`
	Name string			`json:"name" bson:"name"`
	Amount int 			`json:"amount" bson:"amount"`
	Status string		`json:"status" bson:"status"`
}