package main

import (
	"fmt"
	"log"
	"modul1/controller"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", controller.GetAllUsers).Methods("GET")
	router.HandleFunc("/users", controller.InsertNewUser).Methods("POST")
	router.HandleFunc("/users", controller.UpdateUser).Methods("PUT")
	// router.HandleFunc("/users/{id}", controller.DeleteUser).Methods("DELETE")
	router.HandleFunc("/users/{id}", controller.Authenticate(controller.DeleteUser, 2)).Methods("DELETE")
	router.HandleFunc("/users/login", controller.LoginUser).Methods("POST")
	router.HandleFunc("/users/logout", controller.LogoutUser).Methods("POST")

	router.HandleFunc("/products", controller.GetAllProducts).Methods("GET")
	router.HandleFunc("/products", controller.InsertNewProduct).Methods("POST")
	router.HandleFunc("/products", controller.UpdateProduct).Methods("PUT")
	router.HandleFunc("/products/{id}", controller.DeleteProduct).Methods("DELETE")

	router.HandleFunc("/transactions", controller.GetAllTransactions).Methods("GET")
	router.HandleFunc("/transactions", controller.InsertNewTransaction).Methods("POST")
	router.HandleFunc("/transactions", controller.UpdateTransaction).Methods("PUT")
	router.HandleFunc("/transactions/{id}", controller.DeleteTransaction).Methods("DELETE")
	router.HandleFunc("/transactions/detail", controller.GetDetailTransaction).Methods("GET")

	http.Handle("/", router)
	fmt.Println("Connected to port 8080")
	log.Println("Connected to port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))

	//Ini cuma buat mysql bisa diimport
	x := mysql.ErrOldPassword
	fmt.Print(x)

}
