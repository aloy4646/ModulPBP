package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * from transactions"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 150, "Query error")
		return
	}

	var transaction Transaction
	var transactions []Transaction
	for rows.Next() {
		if err := rows.Scan(&transaction.Id, &transaction.Userid, &transaction.Productid, &transaction.Quantity); err != nil {
			log.Println(err.Error())
			errorResponseMessage(w, 170, "Data error")
			return
		} else {
			transactions = append(transactions, transaction)
		}
	}

	var response TransactionsResponse
	response.Status = 200
	response.Message = "Succes"
	response.Data = transactions

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertNewTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 100, "Parse error")
		return
	}

	userId, _ := strconv.Atoi(r.Form.Get("userId"))
	productId, _ := strconv.Atoi(r.Form.Get("productId"))
	quantity, _ := strconv.Atoi(r.Form.Get("quantity"))

	resultQuery, errQuery := db.Exec("INSERT INTO transactions (userId, productId, quantity) VALUES (?,?,?)",
		userId,
		productId,
		quantity,
	)

	productIdNotFoundErr := "Error 1452: Cannot add or update a child row: a foreign key constraint fails (`db_latihan_pbp`.`transactions`, CONSTRAINT `transactions_ibfk_2` FOREIGN KEY (`productId`) REFERENCES `products` (`id`))"

	if errQuery != nil && errQuery.Error() == productIdNotFoundErr {
		if !insertVoidProduct(r.Form.Get("productId"), db) {
			errorResponseMessage(w, 400, "Query error, Insert failed")
			return
		} else {
			resultQuery, errQuery = db.Exec("INSERT INTO transactions (userId, productId, quantity) VALUES (?,?,?)",
				userId,
				productId,
				quantity,
			)
		}
	}

	if errQuery != nil {
		log.Println(errQuery)
		errorResponseMessage(w, 400, "Query error, Insert failed")
		return
	}

	id, _ := resultQuery.LastInsertId()
	var response TransactionResponse
	var transaction Transaction = Transaction{Id: int(id), Userid: userId, Productid: productId, Quantity: quantity}
	response.Status = 200
	response.Message = "Succes"
	response.Data = transaction

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 100, "Parse error")
		return
	}

	id, _ := strconv.Atoi(r.Form.Get("id"))
	userId, _ := strconv.Atoi(r.Form.Get("userId"))
	productId, _ := strconv.Atoi(r.Form.Get("productId"))
	quantity, _ := strconv.Atoi(r.Form.Get("quantity"))

	resultQuery, errQuery := db.Exec("UPDATE transactions SET userId=?, productId=?, quantity=? WHERE id=?",
		userId,
		productId,
		quantity,
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

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	fmt.Println(r)

	vars := mux.Vars(r)
	id := vars["id"]

	resultQuery, err := db.Exec("DELETE FROM transactions WHERE id=?",
		id,
	)

	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 401, "Query error, Delete failed")
	}

	rowsAffected, _ := resultQuery.RowsAffected()
	responseFromRowsAffected(w, rowsAffected)
}

func GetDetailTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	// err := r.ParseForm()
	// if err != nil {
	// 	log.Println(err)
	// 	errorResponseMessage(w, 100, "Parse error")
	// 	return
	// }

	// userId := r.Form.Get("userId")
	// fmt.Print("Id user : ")
	// fmt.Println(userId)

	query := `SELECT a.id, a.quantity, b.*, c.* from transactions a
	JOIN users b on a.userId = b.id
	JOIN products c on a.productId = c.id`

	userId := r.URL.Query()["userId"]
	fmt.Print(userId)
	if userId != nil {
		query += " WHERE userId = " + userId[0]
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 150, "Query error")
		return
	}

	var detailTransaction DetailTransaction
	var user User
	var product Product
	var detailTransactions []DetailTransaction
	for rows.Next() {
		if err := rows.Scan(&detailTransaction.Id, &detailTransaction.Quantity,
			&user.Id, &user.Name, &user.Age, &user.Address,
			&product.Id, &product.Name, &product.Price); err != nil {
			log.Println(err.Error())
			errorResponseMessage(w, 170, "Data error")
			return
		} else {
			detailTransaction.User = user
			detailTransaction.Product = product
			detailTransactions = append(detailTransactions, detailTransaction)
		}
	}

	var response DetailTransactionsResponse
	response.Status = 200
	response.Message = "Succes"
	response.Data = detailTransactions

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func deleteTransactionsByProductId(productId string, db *sql.DB) bool {
	_, err := db.Exec("DELETE FROM transactions WHERE productId=?",
		productId,
	)

	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
