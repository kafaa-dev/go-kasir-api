package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
	"time"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}

func (s *TransactionService) GetReport(from *time.Time, to *time.Time) (*models.ReportResponse, error) {
	totalRevenue, err := s.repo.GetTotalRevenue(from, to)
	if err != nil {
		return nil, err
	}

	totalTransactions, err := s.repo.GetTotalTransactions(from, to)
	if err != nil {
		return nil, err
	}

	bestSellingProduct, err := s.repo.GetBestSellingProduct(from, to)
	if err != nil {
		return nil, err
	}

	return &models.ReportResponse{
		TotalRevenue:       totalRevenue,
		TotalTransactions:  totalTransactions,
		BestSellingProduct: bestSellingProduct,
	}, nil
}
