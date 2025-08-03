package main

import (
	"net/http"
)

func main() {
	homePage := http.FileServer(http.Dir("./"))

	http.Handle("/", homePage)

	http.ListenAndServe(":8080", nil)
}
