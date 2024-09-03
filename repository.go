package main

import (
	"database/sql"
	"errors"
	"fmt"
)

/*
В файле реализованы функции для взаимодействия с базой данных.
Файл содержит методы для получения баланса, добавления, резервирования и списания средств, отката резервирования(отмена), запрос отчета.
*/

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Метод создания нового пользователя со счетом
func (r *Repository) AddNewUser(userID int, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO users (id, balance) VALUES ($1, $2)", userID, amount)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO transactions (user_id, amount, transaction_type) VALUES ($1, $2, 'new_user')", userID, amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Метод зачисления  средтсв на счет пользователя
func (r *Repository) AddDeposit(userID int, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE id = $2", amount, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO transactions (user_id, amount, transaction_type) VALUES ($1, $2, 'add_dep')", userID, amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Метод резервирования средств(reserve)
func (r *Repository) Reserve(userID, serviceID, orderID int, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var balance float64
	row := tx.QueryRow("SELECT balance FROM users WHERE id = $1", userID)
	err = row.Scan(&balance)
	if err != nil {
		return err
	}

	if balance < amount {
		return errors.New("Balance cannot be negative!")
	}

	_, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE id = $2", amount, userID)
	if err != nil {
		return err
	}

	// ставим статус - средства зарезервированы
	_, err = tx.Exec("INSERT INTO transactions (user_id, service_id, order_id, amount, transaction_type) VALUES ($1, $2, $3, $4, 'reserve')", userID, serviceID, orderID, amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Списание средств(reserve_confirm) + добавление  данных в отчет
func (r *Repository) ReserveConfirm(userID, serviceID, orderID int, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// меняем статус - средства списаны
	_, err = tx.Exec("UPDATE transactions SET transaction_type = $1 WHERE user_id = $2 AND service_id = $3 AND order_id = $4", "res_conf", userID, serviceID, orderID)
	if err != nil {
		return err
	}

	// добавление в данных в отчет
	_, err = tx.Exec("INSERT INTO report (service_id, amount) VALUES ($1, $2)", serviceID, amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Откат резервирования средств	(-)
func (r *Repository) ReserveCancel(userID, serviceID, orderID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// меняем  статус на отмену
	_, err = tx.Exec("UPDATE transactions SET transaction_type = $1 WHERE user_id = $2 AND service_id = $3 AND order_id = $4", "res_canc", userID, serviceID, orderID)
	if err != nil {
		return err
	}

	// берем сумму транзакции
	var amount float64
	row := tx.QueryRow("SELECT amount FROM transactions WHERE id = $1 AND service_id = $2 AND order_id = $3", userID, serviceID, orderID)
	err = row.Scan(&amount)
	if err != nil {
		return err
	}

	// возвращаем деньги на счет пользователя
	_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE id = $2", amount, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) GetBalance(userID int) (*User, error) {
	user := &User{}
	err := r.db.QueryRow("SELECT id, balance FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Balance)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetTransactionList() ([]Transactions, error) {
	// Создаем срез для хранения всех транзакций
	transactions := []Transactions{}

	// Выполняем запрос
	selRows, err := r.db.Query("SELECT id, user_id, service_id, order_id, amount, transaction_type FROM transactions")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer selRows.Close()

	// Проходим по всем строкам результата запроса
	for selRows.Next() {
		var tr Transactions
		// Сканируем данные строки в переменные структуры
		err := selRows.Scan(&tr.ID, &tr.UserID, &tr.ServiceID, &tr.OrderID, &tr.Amount, &tr.TransactionType)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		// Добавляем структуру в срез
		transactions = append(transactions, tr)
	}

	// Проверяем на наличие ошибок при итерации
	if err := selRows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return transactions, nil
}

func (r *Repository) GetReportList() ([]Report, error) {
	//repList := &Report{}
	report := []Report{}
	selRows, err := r.db.Query("SELECT * FROM report")
	if err != nil {
		return nil, err
	}
	defer selRows.Close()

	for selRows.Next() {
		var rep Report
		// Сканируем данные строки в переменные структуры
		err := selRows.Scan(&rep.ID, &rep.ServiceID, &rep.Amount)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		// Добавляем структуру в срез
		report = append(report, rep)
	}

	// Проверяем на наличие ошибок при итерации
	if err := selRows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return report, nil
}
