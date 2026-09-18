package main

import (
	"fmt"
	"log"
	"net/http"
	"project/internals/handlers"
	"project/internals/service"
	"project/internals/storage"
)

func main() {

	fp := "info.json"

	s := storage.NewStorage(fp)

	i, err := service.NewService(s)
	if err != nil {
		log.Fatal(err)
	}

	handler := handlers.NewHttpService(i)

	http.HandleFunc("POST /user", handler.HttpCreate)
	http.HandleFunc("GET /user/{id}", handler.HttpGet)
	http.HandleFunc("GET /users", handler.HttpGetAll)
	http.HandleFunc("PUT /user/{id}", handler.HttpUpdate)
	http.HandleFunc("DELETE /user/{id}", handler.HttpDelete)

	fmt.Println("Server running on :8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
