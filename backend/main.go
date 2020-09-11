package main

import (
	"context"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/v12"
	"github.com/opentracing/opentracing-go/log"
	"go-product/backend/web/controllers"
	"go-product/common"
	"go-product/repositories"
	"go-product/services"
)

func main() {
	// 1. create iris instance
	app := iris.New()

	// 2. set error mode on mvc
	app.Logger().SetLevel("debug")

	// 3. Register template
	template := iris.HTML("./backend/web/views", ".html").
		Layout("shared/layout.html").
		Reload(true)
	app.RegisterView(template)

	// 4. set target template
	app.HandleDir("/assets", "./backend/web/assets")

	// Exception redirect page
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问页面出错了"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	// Conenct to database
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()


	// 5. register controller
	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))


	// 6. start server
	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		)

}
