package router

import (
	"ancient-text-backend/internal/handlers"
	"ancient-text-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestID())

	textHandler := handlers.NewTextHandler()
	annotationHandler := handlers.NewAnnotationHandler()
	layoutHandler := handlers.NewLayoutHandler()

	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"service": "ancient-text-backend",
			})
		})

		texts := api.Group("/texts")
		{
			texts.POST("", textHandler.Create)
			texts.GET("", textHandler.List)
			texts.GET("/:id", textHandler.GetByID)
			texts.PUT("/:id", textHandler.Update)
			texts.DELETE("/:id", textHandler.Delete)
		}

		annotations := api.Group("/annotations")
		{
			annotations.POST("", annotationHandler.Create)
			annotations.POST("/batch", annotationHandler.BatchSubmit)
			annotations.GET("/text/:textId", annotationHandler.GetByTextID)
			annotations.GET("/:id", annotationHandler.GetByID)
			annotations.PUT("/:id", annotationHandler.Update)
			annotations.DELETE("/:id", annotationHandler.Delete)
		}

		layouts := api.Group("/layouts")
		{
			layouts.GET("/text/:textId", layoutHandler.GetByTextID)
			layouts.PUT("/text/:textId", layoutHandler.Save)
			layouts.DELETE("/:id", layoutHandler.Delete)
		}
	}

	return r
}
