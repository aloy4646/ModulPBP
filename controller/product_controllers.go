package controller

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * from products"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 150, "Query error")
		return
	}

	var product Product
	var products []Product
	for rows.Next() {
		if err := rows.Scan(&product.Id, &product.Name, &product.Price); err != nil {
			log.Println(err.Error())
			errorResponseMessage(w, 170, "Data error")
			return
		} else {
			products = append(products, product)
		}
	}

	var response ProductsResponse
	response.Status = 200
	response.Message = "Succes"
	response.Data = products

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertNewProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 100, "Parse error")
		return
	}

	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))

	resultQuery, errQuery := db.Exec("INSERT INTO products (name, price) VALUES (?,?)",
		name,
		price,
	)

	if errQuery != nil {
		log.Println(errQuery)
		errorResponseMessage(w, 400, "Query error, Insert failed")
		return
	}

	id, _ := resultQuery.LastInsertId()
	var response ProductResponse
	var product Product = Product{Id: int(id), Name: name, Price: price}
	response.Status = 200
	response.Message = "Succes"
	response.Data = product

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 100, "Parse error")
		return
	}

	id, _ := strconv.Atoi(r.Form.Get("id"))
	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))

	resultQuery, errQuery := db.Exec("UPDATE products SET name=?, price=? WHERE id=?",
		name,
		price,
		id,
	)

	if errQuery != nil {
		log.Println(errQuery)
		errorResponseMessage(w, 400, "Query error, Insert failed")
		return
	}

	rowsAffected, _ := resultQuery.RowsAffected()
	responseFromRowsAffected(w, rowsAffected)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	if !deleteTransactionsByProductId(id, db) {
		errorResponseMessage(w, 401, "Query error, Delete failed")
		return
	}

	resultQuery, err := db.Exec("DELETE FROM products WHERE id=?",
		id,
	)

	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 401, "Query error, Delete failed")
	}

	rowsAffected, _ := resultQuery.RowsAffected()
	responseFromRowsAffected(w, rowsAffected)
}

func insertVoidProduct(producId string, db *sql.DB) bool {
	_, errQuery := db.Exec("INSERT INTO products (id) VALUES (?)",
		producId,
	)

	if errQuery != nil {
		log.Println(errQuery)
		return false
	}
	return true
}
