package datamodels

type Product struct {
	ID           int64  `json:"id" sql:"ID" basic:"ID"`
	ProductName  string `json:"ProductName" sql:"productName" basic:"ProductName"`
	ProductNum   int64  `json:"ProductNum" sql:"productNum" basic:"ProductNum"`
	ProductImage string `json:"ProductImage" sql:"productImage" basic:"ProductImage"`
	ProductUrl   string `json:"ProductUrl" sql:"productUrl" basic:"ProductUrl"`
}
