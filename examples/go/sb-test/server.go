package main

import (
	"net/http"
	"os"
)

func main() {

	api := NewAppRouter()
	http.Handle("/", api)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
