package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT id,name,age,address,email,password from users"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 150, "Query error")
		return
	}

	isValidToken := validateUserToken(r, 2)

	var user User
	var users []User
	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Address, &user.Email, &user.Password); err != nil {
			log.Println(err.Error())
			errorResponseMessage(w, 170, "Data error")
			return
		} else {
			if !isValidToken {
				user.Email = ""
				user.Password = ""
			}
			users = append(users, user)
		}
	}

	var response UsersResponse
	response.Status = 200
	response.Message = "Succes"
	response.Data = users

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertNewUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 100, "Parse error")
		return
	}

	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	resultQuery, errQuery := db.Exec("INSERT INTO users (name, age, address,email,password) VALUES (?,?,?)",
		name,
		age,
		address,
		email,
		password,
	)

	if errQuery != nil {
		log.Println(errQuery)
		errorResponseMessage(w, 400, "Query error, Insert failed")
		return
	}

	id, _ := resultQuery.LastInsertId()
	var response UserResponse
	var user User = User{Id: int(id), Name: name, Age: age, Address: address}
	response.Status = 200
	response.Message = "Succes"
	response.Data = user

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
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
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")

	resultQuery, errQuery := db.Exec("UPDATE users SET name=?, age=?, address=? WHERE id=?",
		name,
		age,
		address,
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

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	resultQuery, err := db.Exec("DELETE FROM users WHERE id=?",
		id,
	)

	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 401, "Query error, Delete failed")
	}

	rowsAffected, _ := resultQuery.RowsAffected()
	responseFromRowsAffected(w, rowsAffected)
}

func responseFromRowsAffected(w http.ResponseWriter, rowsAffected int64) {
	if rowsAffected > 0 {
		successResponseMessage(w)
	} else {
		errorResponseMessage(w, 407, "Failed, 0 rows affected")
	}
}

func successResponseMessage(w http.ResponseWriter) {
	var response SuccessResponse
	response.Status = 200
	response.Message = "Success"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func errorResponseMessage(w http.ResponseWriter, status int, message string) {
	var response ErrorResponse
	response.Status = status
	response.Message = message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 100, "Parse error")
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	var user User
	if err := db.QueryRow("SELECT id, name, age, address from users where email = ? AND password = ?",
		email, password).Scan(&user.Id, &user.Name, &user.Age, &user.Address); err != nil {
		log.Println(err.Error())
		errorResponseMessage(w, 170, "Login gagal")
		return
	}

	//hanya untuk percobaan, karena DB belum diubah
	user.UserType, _ = strconv.Atoi(r.Form.Get("user type"))
	generateToken(w, user.Id, user.Name, user.UserType)

	var response UserResponse
	response.Status = 200
	response.Message = "Login Success"
	response.Data = user

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	resetUserToken(w)
	successResponseMessage(w)
}

func sendUnAuthorizedResponse(w http.ResponseWriter) {
	var response ErrorResponse
	response.Status = 401
	response.Message = "UnAuthorized Access"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
