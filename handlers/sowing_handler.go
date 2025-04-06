package handlers

import (
	"crop-tracker/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateSowing(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sowing models.Sowing
		if err := c.BindJSON(&sowing); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate required fields
		if sowing.FieldID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "field_id is required and must be positive"})
			return
		}
		if sowing.Crop == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "crop is required"})
			return
		}
		if sowing.SowedAt == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "sowed_at is required"})
			return
		}

		// Check if field exists
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM fields WHERE id = ?)", sowing.FieldID).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "field does not exist"})
			return
		}

		res, err := db.Exec("INSERT INTO sowings(field_id, crop, sowed_at) VALUES(?, ?, ?)",
			sowing.FieldID, sowing.Crop, sowing.SowedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		id, _ := res.LastInsertId()
		sowing.ID = int(id)
		c.JSON(http.StatusOK, sowing)
	}
}

func ListSowings(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, field_id, crop, sowed_at FROM sowings")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var sowings []models.Sowing
		for rows.Next() {
			var s models.Sowing
			if err := rows.Scan(&s.ID, &s.FieldID, &s.Crop, &s.SowedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			sowings = append(sowings, s)
		}

		c.JSON(http.StatusOK, sowings)
	}
}
