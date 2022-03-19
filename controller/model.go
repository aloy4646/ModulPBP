package controller

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
	Email    string `json:"email"`
	Password string `json:"password"`
	UserType int    `json:"user_type"`
}

type Product struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Transaction struct {
	Id        int `json:"id"`
	Userid    int `json:"User ID"`
	Productid int `json:"Product ID"`
	Quantity  int `json:"Quantity"`
}

type DetailTransaction struct {
	Id       int     `json:"id"`
	User     User    `json:"User"`
	Product  Product `json:"Product"`
	Quantity int     `json:"Quantity"`
}

type UsersResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []User `json:"users"`
}

type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    User   `json:"user"`
}

type ProductsResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    []Product `json:"Products"`
}

type ProductResponse struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Data    Product `json:"Products"`
}

type TransactionsResponse struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    []Transaction `json:"Transaction"`
}

type TransactionResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    Transaction `json:"Transaction"`
}

type DetailTransactionsResponse struct {
	Status  int                 `json:"status"`
	Message string              `json:"message"`
	Data    []DetailTransaction `json:"Detail Transaction"`
}

type SuccessResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
