package handlers

import (
	"crop-tracker/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddHarvest(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var harvest models.Harvest
		if err := c.BindJSON(&harvest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate required fields
		if harvest.FieldID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "field_id is required and must be positive"})
			return
		}
		if harvest.Crop == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "crop is required"})
			return
		}
		if harvest.YieldPerHa < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "yield_t_per_ha must be non-negative"})
			return
		}

		// Check if field exists
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM fields WHERE id = ?)", harvest.FieldID).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "field does not exist"})
			return
		}

		res, err := db.Exec("INSERT INTO harvests(field_id, crop, yield_t_per_ha) VALUES(?, ?, ?)",
			harvest.FieldID, harvest.Crop, harvest.YieldPerHa)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		id, _ := res.LastInsertId()
		harvest.ID = int(id)
		c.JSON(http.StatusOK, harvest)
	}
}

func ListHarvests(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, field_id, crop, yield_t_per_ha FROM harvests")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var harvests []models.Harvest
		for rows.Next() {
			var h models.Harvest
			if err := rows.Scan(&h.ID, &h.FieldID, &h.Crop, &h.YieldPerHa); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			harvests = append(harvests, h)
		}

		c.JSON(http.StatusOK, harvests)
	}
}
