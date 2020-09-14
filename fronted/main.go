package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/opentracing/opentracing-go/log"
	"go-product/common"
	"go-product/fronted/web/controllers"
	"go-product/repositories"
	"go-product/services"
	"time"
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

	// init session
	sess := sessions.New(sessions.Config{
		Cookie:  "AdminCookie",
		Expires: 600 * time.Minute,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user := repositories.NewUserRepository("user", db)
	userService := services.NewService(user)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService, ctx, sess.Start)
	userPro.Handle(new(controllers.UserController))

	app.Run(
		iris.Addr("0.0.0.0:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}
