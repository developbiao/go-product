package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

/*
编写 HTTPS 服务器
HTTPS = HTTP + Secure(安全)

RSA 进行加密
SHA 进行验证
密钥和证书

生成密钥文件
openssl genrsa -out server.key 2048

生成证书文件
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650

*/

func httpsServer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path:", r.URL.Path)
	fmt.Println("Url:", r.URL)
	fmt.Println("Host:", r.Host)
	fmt.Println("Header:", r.Header)
	fmt.Println("Method:", r.Method)
	fmt.Println("Proto:", r.Proto)
	fmt.Println("UserAgent:", r.UserAgent())

	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	fmt.Println("完整的请求路径:", strings.Join([]string{scheme, r.Host, r.RequestURI}, ""))
	fmt.Fprintf(w, "Hello Go Web")
}

func main() {
	http.HandleFunc("/", httpsServer)
	fmt.Println("HTTPS 服务器已经启动，请在浏览器地址栏中输入 https://localhost:433/")
	err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
