package handlers

import (
	"bytes"
	"crop-tracker/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE fields (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			area_ha REAL NOT NULL,
			region TEXT NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	return db
}

func TestCreateField(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	defer db.Close()

	tests := []struct {
		name       string
		field      models.Field
		wantStatus int
		wantError  bool
	}{
		{
			name: "Valid field",
			field: models.Field{
				Name:   "Test Field",
				AreaHa: 100.5,
				Region: "Test Region",
			},
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name: "Invalid field - empty name",
			field: models.Field{
				AreaHa: 100.5,
				Region: "Test Region",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create request body
			body, _ := json.Marshal(tt.field)
			c.Request = httptest.NewRequest("POST", "/fields", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call handler
			CreateField(db)(c)

			assert.Equal(t, tt.wantStatus, w.Code)

			if !tt.wantError {
				var response models.Field
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotZero(t, response.ID)
				assert.Equal(t, tt.field.Name, response.Name)
				assert.Equal(t, tt.field.AreaHa, response.AreaHa)
				assert.Equal(t, tt.field.Region, response.Region)
			}
		})
	}
}

func TestListFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	testFields := []models.Field{
		{Name: "Field 1", AreaHa: 100.5, Region: "Region 1"},
		{Name: "Field 2", AreaHa: 200.7, Region: "Region 2"},
	}

	for _, field := range testFields {
		_, err := db.Exec("INSERT INTO fields(name, area_ha, region) VALUES(?, ?, ?)",
			field.Name, field.AreaHa, field.Region)
		if err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}
	}

	// Test listing fields
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/fields", nil)

	ListFields(db)(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Field
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, len(testFields))

	// Verify field contents
	for i, field := range response {
		assert.NotZero(t, field.ID)
		assert.Equal(t, testFields[i].Name, field.Name)
		assert.Equal(t, testFields[i].AreaHa, field.AreaHa)
		assert.Equal(t, testFields[i].Region, field.Region)
	}
}
