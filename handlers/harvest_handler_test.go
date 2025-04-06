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

func setupTestDBWithSowings(t *testing.T) *sql.DB {
	db := setupTestDBWithFields(t)

	// Create harvests table
	_, err := db.Exec(`
		CREATE TABLE harvests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			field_id INTEGER NOT NULL,
			crop TEXT NOT NULL,
			yield_t_per_ha REAL NOT NULL,
			FOREIGN KEY(field_id) REFERENCES fields(id)
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create harvests table: %v", err)
	}

	return db
}

func TestAddHarvest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDBWithSowings(t)
	defer db.Close()

	tests := []struct {
		name       string
		harvest    models.Harvest
		wantStatus int
		wantError  bool
	}{
		{
			name: "Valid harvest",
			harvest: models.Harvest{
				FieldID:    1,
				Crop:       "Wheat",
				YieldPerHa: 4.5,
			},
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name: "Invalid field ID",
			harvest: models.Harvest{
				FieldID:    999,
				Crop:       "Wheat",
				YieldPerHa: 4.5,
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name: "Missing crop",
			harvest: models.Harvest{
				FieldID:    1,
				YieldPerHa: 4.5,
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name: "Negative yield",
			harvest: models.Harvest{
				FieldID:    1,
				Crop:       "Wheat",
				YieldPerHa: -1.0,
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.harvest)
			c.Request = httptest.NewRequest("POST", "/harvest", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			AddHarvest(db)(c)

			assert.Equal(t, tt.wantStatus, w.Code)

			if !tt.wantError {
				var response models.Harvest
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotZero(t, response.ID)
				assert.Equal(t, tt.harvest.FieldID, response.FieldID)
				assert.Equal(t, tt.harvest.Crop, response.Crop)
				assert.Equal(t, tt.harvest.YieldPerHa, response.YieldPerHa)
			}
		})
	}
}

func TestListHarvests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDBWithSowings(t)
	defer db.Close()

	// Insert test harvests
	testHarvests := []models.Harvest{
		{FieldID: 1, Crop: "Wheat", YieldPerHa: 4.5},
		{FieldID: 1, Crop: "Corn", YieldPerHa: 6.7},
	}

	for _, harvest := range testHarvests {
		_, err := db.Exec("INSERT INTO harvests(field_id, crop, yield_t_per_ha) VALUES(?, ?, ?)",
			harvest.FieldID, harvest.Crop, harvest.YieldPerHa)
		if err != nil {
			t.Fatalf("Failed to insert test harvest: %v", err)
		}
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/harvest", nil)

	ListHarvests(db)(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Harvest
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, len(testHarvests))

	for i, harvest := range response {
		assert.NotZero(t, harvest.ID)
		assert.Equal(t, testHarvests[i].FieldID, harvest.FieldID)
		assert.Equal(t, testHarvests[i].Crop, harvest.Crop)
		assert.Equal(t, testHarvests[i].YieldPerHa, harvest.YieldPerHa)
	}
}
