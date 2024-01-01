package main

import (
	"fmt"
	"net/http"// 公式が提供するパッケージ
)

type myHandler struct{}

// myHandlerに対して、HTTPハンドラーを提供する
func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}


func main()  {
	fmt.Println("Hello World")

	handler := &myHandler{}
	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)
}