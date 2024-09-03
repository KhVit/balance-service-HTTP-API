package main

/* обработка HTTP-запросов, используя методы из сервиса */

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Handlers struct {
	service *Service
}

func NewHandlers(service *Service) *Handlers {
	return &Handlers{service: service}
}

// Обработчик для создания нового пользователя со счетом
func (h *Handlers) HandleAddNewUser(w http.ResponseWriter, r *http.Request) {
	//var req struct {
	//	UserID int     `json:"user_id"`
	//	Amount float64 `json:"balance"`
	//}
	fmt.Println("POST /new-user - HandleAddNewUser")

	user := &User{}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.AddNewUser(user.ID, user.Balance); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to deposit funds", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// Обработчик для зачисления средств на баланс
func (h *Handlers) HandleAddDeposit(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT /add-deposit - HandleAddDeposit")

	user := &User{}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.AddDeposit(user.ID, user.Balance); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to deposit funds", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// Обработчик для резервирования средств
func (h *Handlers) HandleReserve(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT /reserve - HandleReserve")

	tran := Transactions{}

	if err := json.NewDecoder(r.Body).Decode(&tran); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Reserve(tran.UserID, tran.ServiceID, tran.OrderID, tran.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// Обработчик для подтверждения списания средств
func (h *Handlers) HandleReserveConfirm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT /reserve-confirm - HandleReserveConfirm")

	tran := Transactions{}

	if err := json.NewDecoder(r.Body).Decode(&tran); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.ReserveConfirm(tran.UserID, tran.ServiceID, tran.OrderID, tran.Amount); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// Обработчик для отмены списания(резервирования) средств
func (h *Handlers) HandleReserveCancel(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT /reserve_cancel - HandleReserveCancel")

	tran := Transactions{}

	if err := json.NewDecoder(r.Body).Decode(&tran); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.ReserveCancel(tran.UserID, tran.ServiceID, tran.OrderID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// Обработчик для получения баланса пользователя
func (h *Handlers) HandleGetUserBalance(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /user-balance/{id} - HandleGetUserBalance")

	userIDStr := strings.TrimPrefix(r.URL.Path, "/userBalance/")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	//userID := r.URL.Path("/balance")

	user, err := h.service.GetBalance(userID)
	if err != nil {
		http.Error(w, "Failed to get balance", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Обработчик для получения баланса пользователя
func (h *Handlers) HandleGetTransactionList(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("GET /tran-list - HandleGetTransactionList")

	tranList, err := h.service.GetTransactionList()
	if err != nil {
		http.Error(w, "Failed to get Transaction list!", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tranList)
}

// Обработчик для получения отчета
func (h *Handlers) HandleGetReport(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("GET /report - HandleGetReport")

	repList, err := h.service.GetReportList()
	if err != nil {
		http.Error(w, "Failed to get report!", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repList)
}
