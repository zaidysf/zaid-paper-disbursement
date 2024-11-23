package services

import (
	"database/sql"
	"errors"
	"fmt"
)

type DisbursementService struct {
	db *sql.DB
}

func NewDisbursementService(db *sql.DB) *DisbursementService {
	return &DisbursementService{db: db}
}

type DisbursementRequest struct {
	UserID int64   `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func (s *DisbursementService) ProcessDisbursement(req *DisbursementRequest) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Check user balance
	var currentBalance float64
	err = tx.QueryRow(`
		SELECT balance 
		FROM user_balances 
		WHERE user_id = $1 
		FOR UPDATE
	`, req.UserID).Scan(&currentBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user balance not found")
		}
		return fmt.Errorf("error checking balance: %v", err)
	}

	// Validate sufficient balance
	if currentBalance < req.Amount {
		return errors.New("insufficient balance")
	}

	// Update balance
	_, err = tx.Exec(`
		UPDATE user_balances 
		SET balance = balance - $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $2
	`, req.Amount, req.UserID)
	if err != nil {
		return fmt.Errorf("error updating balance: %v", err)
	}

	// Record transaction
	_, err = tx.Exec(`
		INSERT INTO transactions (user_id, amount, type)
		VALUES ($1, $2, $3)
	`, req.UserID, req.Amount, "disbursement")
	if err != nil {
		return fmt.Errorf("error recording transaction: %v", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}