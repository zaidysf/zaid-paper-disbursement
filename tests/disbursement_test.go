package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"zaid-paper-disbursement/api/handlers"
	"zaid-paper-disbursement/config"
	"zaid-paper-disbursement/internal/services"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	// Setup
	var err error
	testDB, err = config.InitTestDB()
	if err != nil {
		panic(err)
	}
	defer testDB.Close()

	// Setup test database schema
	if err := config.SetupTestDB(testDB); err != nil {
		panic(err)
	}

	// Run tests
	m.Run()
}

func setupTestRouter(db *sql.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	disbursementService := services.NewDisbursementService(db)
	disbursementHandler := handlers.NewDisbursementHandler(disbursementService)

	router.POST("/api/v1/disbursement", disbursementHandler.ProcessDisbursement)
	return router
}

func setupTestData(db *sql.DB) (int64, error) {
	// Clean previous test data
	if err := config.CleanTestDB(db); err != nil {
		return 0, err
	}

	// Insert test user
	var userID int64
	err := db.QueryRow(`
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id
	`, "testuser", "test@example.com").Scan(&userID)
	if err != nil {
		return 0, err
	}

	// Set initial balance
	_, err = db.Exec(`
		INSERT INTO user_balances (user_id, balance)
		VALUES ($1, $2)
	`, userID, 1000.00)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func TestProcessDisbursement(t *testing.T) {
	userID, err := setupTestData(testDB)
	if err != nil {
		t.Fatalf("Failed to setup test data: %v", err)
	}

	router := setupTestRouter(testDB)

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Valid disbursement",
			payload: map[string]interface{}{
				"user_id": userID,
				"amount":  100.00,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid amount",
			payload: map[string]interface{}{
				"user_id": userID,
				"amount":  -100.00,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Insufficient balance",
			payload: map[string]interface{}{
				"user_id": userID,
				"amount":  10000.00,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/api/v1/disbursement", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}