package main

import (
	"log"
	"net/http"
	"time"
)

func helloHandler(wr http.ResponseWriter, r *http.Request) {
	wr.Write([]byte("welcome!"))
}

// middleware
func timeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()

		// next Handler
		next.ServeHTTP(wr, r)

		timeElapsed := time.Since(timeStart)
		log.Println(timeElapsed)
	})
}

func main() {
	http.Handle("/", timeMiddleware(http.HandlerFunc(helloHandler)))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
