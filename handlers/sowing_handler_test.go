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

func setupTestDBWithFields(t *testing.T) *sql.DB {
	db := setupTestDB(t)

	// Create sowings table
	_, err := db.Exec(`
		CREATE TABLE sowings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			field_id INTEGER NOT NULL,
			crop TEXT NOT NULL,
			sowed_at TEXT NOT NULL,
			FOREIGN KEY(field_id) REFERENCES fields(id)
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create sowings table: %v", err)
	}

	// Insert a test field
	_, err = db.Exec("INSERT INTO fields(name, area_ha, region) VALUES(?, ?, ?)",
		"Test Field", 100.5, "Test Region")
	if err != nil {
		t.Fatalf("Failed to insert test field: %v", err)
	}

	return db
}

func TestCreateSowing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDBWithFields(t)
	defer db.Close()

	tests := []struct {
		name       string
		sowing     models.Sowing
		wantStatus int
		wantError  bool
	}{
		{
			name: "Valid sowing",
			sowing: models.Sowing{
				FieldID: 1,
				Crop:    "Wheat",
				SowedAt: "2025-04-06",
			},
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name: "Invalid field ID",
			sowing: models.Sowing{
				FieldID: 999,
				Crop:    "Wheat",
				SowedAt: "2025-04-06",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name: "Missing crop",
			sowing: models.Sowing{
				FieldID: 1,
				SowedAt: "2025-04-06",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.sowing)
			c.Request = httptest.NewRequest("POST", "/sowings", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			CreateSowing(db)(c)

			assert.Equal(t, tt.wantStatus, w.Code)

			if !tt.wantError {
				var response models.Sowing
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotZero(t, response.ID)
				assert.Equal(t, tt.sowing.FieldID, response.FieldID)
				assert.Equal(t, tt.sowing.Crop, response.Crop)
				assert.Equal(t, tt.sowing.SowedAt, response.SowedAt)
			}
		})
	}
}

func TestListSowings(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDBWithFields(t)
	defer db.Close()

	// Insert test sowings
	testSowings := []models.Sowing{
		{FieldID: 1, Crop: "Wheat", SowedAt: "2025-04-06"},
		{FieldID: 1, Crop: "Corn", SowedAt: "2025-04-07"},
	}

	for _, sowing := range testSowings {
		_, err := db.Exec("INSERT INTO sowings(field_id, crop, sowed_at) VALUES(?, ?, ?)",
			sowing.FieldID, sowing.Crop, sowing.SowedAt)
		if err != nil {
			t.Fatalf("Failed to insert test sowing: %v", err)
		}
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/sowings", nil)

	ListSowings(db)(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Sowing
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, len(testSowings))

	for i, sowing := range response {
		assert.NotZero(t, sowing.ID)
		assert.Equal(t, testSowings[i].FieldID, sowing.FieldID)
		assert.Equal(t, testSowings[i].Crop, sowing.Crop)
		assert.Equal(t, testSowings[i].SowedAt, sowing.SowedAt)
	}
}
