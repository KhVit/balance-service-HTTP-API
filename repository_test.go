// repository_test.go

package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTransactionList(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Определяем, какие данные должна вернуть база данных
	rows := sqlmock.NewRows([]string{"id", "user_id", "service_id", "order_id", "amount", "transaction_type"}).
		AddRow(1, 123, 1, 1, 100.0, "debit").
		AddRow(2, 124, 2, 2, 200.0, "credit")

	// Ожидаем выполнение запроса и возвращение результатов
	mock.ExpectQuery("SELECT id, user_id, service_id, order_id, amount, transaction_type FROM transactions").WillReturnRows(rows)

	// Создаем репозиторий с моком базы данных
	repo := repository.NewRepository(db)

	// Вызываем тестируемую функцию
	transactions, err := repo.GetTransactionList()

	// Проверяем, что ошибок нет
	assert.NoError(t, err)
	// Проверяем, что мы получили правильное количество транзакций
	assert.Len(t, transactions, 2)

	// Проверяем значения первой транзакции
	assert.Equal(t, 1, transactions[0].ID)
	assert.Equal(t, 123, transactions[0].UserID)
	assert.Equal(t, "debit", transactions[0].TransactionType)

	// Проверяем значения второй транзакции
	assert.Equal(t, 2, transactions[1].ID)
	assert.Equal(t, 124, transactions[1].UserID)
	assert.Equal(t, "credit", transactions[1].TransactionType)

	// Проверяем, что все ожидания были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}
