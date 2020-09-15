package middleware

import "github.com/kataras/iris/v12"

func AuthConProduct(ctx iris.Context) {

	uid := ctx.GetCookie("uid")
	if uid == "" {
		ctx.Application().Logger().Debug("Must login before!")
		ctx.Redirect("/user/login")
		return
	}
	ctx.Application().Logger().Debug("Aready login")
	ctx.Next()
}
