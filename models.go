package main

// юзеры
type User struct {
	ID      int     `json:"id"`
	Balance float64 `json:"balance"`
}

// транзакции
type Transactions struct {
	ID              int     `json:"id"`
	UserID          int     `json:"user_id"`
	ServiceID       int     `json:"service_id,omitempty"`
	OrderID         int     `json:"order_id,omitempty"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"` // "deposit", "reserve", "reserve_confirm", "reserve_cancel"
}

// отчет для бухгалтерии
type Report struct {
	ID        int     `json:"id"`
	ServiceID int     `json:"service_id"`
	Amount    float64 `json:"amount"`
}
