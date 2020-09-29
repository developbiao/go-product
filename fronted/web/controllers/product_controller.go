package controllers

import (
	"encoding/json"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"go-product/datamodels"
	"go-product/rabbitmq"
	"go-product/services"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	RabbitMQ       *rabbitmq.RabbitMQ
	Session        *sessions.Session
}

var (
	// Generate html path
	htmlOutPath = "./fronted/web/htmlProductShow/"
	// Static template path
	templatePath = "./fronted/web/views/template/"
)

// Generate product static page
func (p *ProductController) GetGenerateHtml() {
	// Get product id
	productString := p.Ctx.URLParam("productID")
	productID, err := strconv.Atoi(productString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	// Get template
	contensTmp, err := template.ParseFiles(filepath.Join(templatePath, "product.html"))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	// Get output html path
	fileName := filepath.Join(htmlOutPath, "htmlProduct.html")

	// Get render data
	product, err := p.ProductService.GetProductByID(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	// Generate status file
	generateStaticHtml(p.Ctx, contensTmp, fileName, product)

}

func generateStaticHtml(ctx iris.Context, template *template.Template,
	fileName string, product *datamodels.Product) {
	// 1. check file exists?
	if exist(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			ctx.Application().Logger().Error(err)
		}
	}

	//  2. Generate status file
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		ctx.Application().Logger().Error(err)
	}
	defer file.Close()
	template.Execute(file, &product)

}

func exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func (p *ProductController) GetDetail() mvc.View {
	productString := p.Ctx.URLParam("productID")
	productID, err := strconv.Atoi(productString)
	if err != nil {
		productID = 1
	}
	product, err := p.ProductService.GetProductByID(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
	}

	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

func (p *ProductController) GetOrder() []byte {
	productString := p.Ctx.URLParam("productID")
	userString := p.Ctx.GetCookie("uid")
	productID, err := strconv.ParseInt(productString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	// Optimization using rabbitmq
	userID, err := strconv.ParseInt(userString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	// Create message body
	message := datamodels.NewMessage(userID, productID)

	// Type convert to json string
	byteMessage, err := json.Marshal(message)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	err = p.RabbitMQ.PublishSimple(string(byteMessage))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return []byte("true")

	//product, err := p.ProductService.GetProductByID(int64(productID))
	//if err != nil {
	//	p.Ctx.Application().Logger().Debug(err)
	//}
	//
	//var orderID int64
	//showMessage := "抢购失败!" // default message
	//
	//// judge product number
	//if product.ProductNum > 0 {
	//	// reduce product amount
	//	product.ProductNum -= 1
	//	err := p.ProductService.UpdateProduct(product)
	//	if err != nil {
	//		p.Ctx.Application().Logger().Debug(err)
	//	}
	//
	//	// create order
	//	userID, err := strconv.Atoi(userString)
	//	if err != nil {
	//		p.Ctx.Application().Logger().Debug(err)
	//	}
	//
	//	order := &datamodels.Order{
	//		UserId:      int64(userID),
	//		ProductId:   int64(productID),
	//		OrderStatus: datamodels.OrderSuccess,
	//	}
	//
	//	// insert order
	//	orderID, err = p.OrderService.InsertOrder(order)
	//	if err != nil {
	//		p.Ctx.Application().Logger().Debug(err)
	//	} else {
	//		showMessage = "抢购成功!"
	//	}
	//}

	//return mvc.View{
	//	Layout: "shared/productLayout.html",
	//	Name:   "product/result.html",
	//	Data: iris.Map{
	//		"orderID":     orderID,
	//		"showMessage": showMessage,
	//	},
	//}

}
