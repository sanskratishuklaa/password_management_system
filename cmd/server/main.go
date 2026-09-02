package main

import (
	"log"
	"net/http"
	"password_manager/config"
	"password_manager/internal/handler"
)

func main(){
	db, err:=config.ConnectDatabase()
	if err!=nil{
		log.Fatal(err)
	}
	defer db.Close(nil)
	log.Println("Database Connected Successfully")
	http.HandleFunc("/api/v1/health", handler.Health)

	log.Println("Server running on http://localhost:8080")

	err =http.ListenAndServe(":8080",nil)
	if err!=nil{
		log.Fatal(err)
	}
}
