package handlers

import (
	"crop-tracker/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateField(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var field models.Field
		if err := c.BindJSON(&field); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate required fields
		if field.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		if field.AreaHa <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "area_ha must be positive"})
			return
		}
		if field.Region == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "region is required"})
			return
		}

		res, err := db.Exec("INSERT INTO fields(name, area_ha, region) VALUES(?, ?, ?)",
			field.Name, field.AreaHa, field.Region)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		id, _ := res.LastInsertId()
		field.ID = int(id)
		c.JSON(http.StatusOK, field)
	}
}

func ListFields(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name, area_ha, region FROM fields")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var fields []models.Field
		for rows.Next() {
			var field models.Field
			if err := rows.Scan(&field.ID, &field.Name, &field.AreaHa, &field.Region); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			fields = append(fields, field)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, fields)
	}
}
