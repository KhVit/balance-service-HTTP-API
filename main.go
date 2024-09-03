package main

/*Тестовое задание avito.tech 2022*/

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=postgres password=12345 dbname=balance_service sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// Инициализация репозитория, сервиса и обработчиков
	repo := NewRepository(db)
	service := NewService(repo)
	handlers := NewHandlers(service)

	//http.HandleFunc("/balance/deposit", handlers.HandleDeposit)
	//http.HandleFunc("/balance/reserve", handlers.HandleReserve)
	//http.HandleFunc("/balance/confirm", handlers.HandleConfirm)
	//http.HandleFunc("/balance/", handlers.HandleGetBalance)
	//
	//fmt.Println("Starting server on :8080...")
	//log.Fatal(http.ListenAndServe(":8080", nil))

	mux := http.NewServeMux()
	mux.HandleFunc("/newUser", handlers.HandleAddNewUser)
	mux.HandleFunc("/addDeposit", handlers.HandleAddDeposit)
	mux.HandleFunc("/reserve", handlers.HandleReserve)
	mux.HandleFunc("/reserveConfirm", handlers.HandleReserveConfirm)
	mux.HandleFunc("/reserveCancel", handlers.HandleReserveCancel)
	mux.HandleFunc("/userBalance/{user_id}", handlers.HandleGetUserBalance)
	mux.HandleFunc("/tranList", handlers.HandleGetTransactionList)
	mux.HandleFunc("/reportList", handlers.HandleGetReport)

	fmt.Println("Starting server on :8080...")
	//log.Fatal(http.ListenAndServe(":8080", nil))
	errLS := http.ListenAndServe(":8080", mux)
	if errLS != nil {
		fmt.Println(errLS)
	}
}
