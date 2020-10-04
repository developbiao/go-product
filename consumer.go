package main

import (
	"fmt"
	"go-product/common"
	"go-product/rabbitmq"
	"go-product/repositories"
	"go-product/services"
)

func main() {
	// Product consumer
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}

	// Create product database operator instance
	product := repositories.NewProductManager("product", db)
	// Create product service
	productService := services.NewProductService(product)
	// Create order database instance
	order := repositories.NewOrderMangerRepository("order", db)

	// create order service
	orderService := services.NewOrderService(order)

	rabbitmqConsumeSimple := rabbitmq.NewRabbitMQSimple("golang")
	rabbitmqConsumeSimple.ConsumeSimple(orderService, productService)

}
