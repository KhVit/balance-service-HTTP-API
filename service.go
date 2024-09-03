package main

/*
В файле определена бизнес-логика, которая будет использовать методы репозитория. Это уровень абстракции, который связывает репозиторий с обработчиками HTTP-запросов
*/

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddNewUser(userID int, amount float64) error {
	return s.repo.AddNewUser(userID, amount)
}

func (s *Service) AddDeposit(userID int, amount float64) error {
	return s.repo.AddDeposit(userID, amount)
}

func (s *Service) Reserve(userID, serviceID, orderID int, amount float64) error {
	return s.repo.Reserve(userID, serviceID, orderID, amount)
}

func (s *Service) ReserveConfirm(userID, serviceID, orderID int, amount float64) error {
	return s.repo.ReserveConfirm(userID, serviceID, orderID, amount)
}

func (s *Service) ReserveCancel(userID, serviceID, orderID int) error {
	return s.repo.ReserveCancel(userID, serviceID, orderID)
}

func (s *Service) GetBalance(userID int) (*User, error) {
	return s.repo.GetBalance(userID)
}

func (s *Service) GetTransactionList() ([]Transactions, error) {
	return s.repo.GetTransactionList()
}

func (s *Service) GetReportList() ([]Report, error) {
	return s.repo.GetReportList()
}
