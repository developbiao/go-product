package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/opentracing/opentracing-go/log"
	"go-product/common"
	"go-product/fronted/middleware"
	"go-product/fronted/web/controllers"
	"go-product/repositories"
	"go-product/services"
)

func main() {

	// create iris instance
	app := iris.New()

	// set log level
	app.Logger().SetLevel("debug")

	// register template
	template := iris.HTML("./fronted/web/views", ".html").Layout("shared/layout.html")
	app.RegisterView(template)

	// set template
	app.HandleDir("/public", "./fronted/web/public")
	app.HandleDir("/html", "./fronted/web/htmlProductShow")

	//  Exception redirect page
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "Visited page error!"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	//  connect to database
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Error(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// user controller register
	user := repositories.NewUserRepository("user", db)
	userService := services.NewService(user)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService, ctx)
	userPro.Handle(new(controllers.UserController))

	// product controller register
	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderMangerRepository("order", db)
	orderService := services.NewOrderService(order)
	proProduct := app.Party("/product")
	pro := mvc.New(proProduct)
	proProduct.Use(middleware.AuthConProduct)
	pro.Register(productService, orderService)
	pro.Handle(new(controllers.ProductController))

	app.Run(
		iris.Addr("0.0.0.0:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}
