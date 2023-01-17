package main

import (
	"fmt"
	"log"
	router "mongoapi/routers"
	"net/http"
)

func main() {
	fmt.Println("Server Started at PORT 4000")
	r := router.Router()
	log.Fatal(http.ListenAndServe(":4000", r))
}