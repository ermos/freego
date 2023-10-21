package main

import (
	"fmt"
	"net/http"
)

func main() {
	go addRouter("hello, app1 !", "8080")
	go addRouter("hello, app2 !", "8081")
	go addRouter("hello, all !", "8082")
	addRouter("hello, world !", "8083")
}

func addRouter(word string, port string) {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, word)
	})
	http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}
