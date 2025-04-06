package routers

import (
	"crop-tracker/handlers"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	r.POST("/fields", handlers.CreateField(db))
	r.GET("/fields", handlers.ListFields(db))

	r.POST("/sowings", handlers.CreateSowing(db))
	r.GET("/sowings", handlers.ListSowings(db))

	r.POST("/harvest", handlers.AddHarvest(db))
	r.GET("/harvest", handlers.ListHarvests(db))

	return r
}
