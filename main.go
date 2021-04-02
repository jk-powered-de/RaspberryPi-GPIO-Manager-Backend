package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"raspberrypi-gpio-manager-backend/db"
	"raspberrypi-gpio-manager-backend/handler"
	"raspberrypi-gpio-manager-backend/model"
)

func main() {
	dbConfig := db.LoadDatabaseConfig()
	db.ConnectDatabase(dbConfig)

	err := db.Connection.AutoMigrate(model.Job{})
	if err != nil {
		panic(err)
	}

	err = db.Connection.AutoMigrate(model.NamedPin{})
	if err != nil {
		panic(err)
	}

	go handler.StartInterval()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(cors.Default())

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code":    "404",
			"message": "not_found",
		})
	})

	r.HandleMethodNotAllowed = true
	r.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{
			"code":    "405",
			"message": "method_not_allowed",
		})
	})

	r.GET("/", handler.ServiceInfo)

	v1 := r.Group("/v1")
	{
		v1.POST("/named-pin/:id/action/turnon", handler.TurnOnNamedPinByID)
		v1.POST("/named-pin/:id/action/turnoff", handler.TurnOffNamedPinByID)

		v1.GET("/named-pins", handler.FindAllNamedPins)
		v1.POST("/named-pins", handler.CreateNamedPin)
		v1.PATCH("/named-pins/:id", handler.UpdateNamedPinByID)
		v1.DELETE("/named-pins/:id", handler.DeleteNamedPinByID)

		v1.GET("/jobs", handler.FindAllJobs)
		v1.GET("/jobs/undone", handler.FindAllJobsUndone)
		v1.GET("/jobs/done", handler.FindAllJobsDone)
		v1.POST("/jobs", handler.CreateJob)
		v1.DELETE("/jobs/:id", handler.DeleteJobByID)
		v1.DELETE("/jobs/:id/named-pin", handler.DeleteJobByNamedGpioPinId)
	}

	err = r.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

// Ideas: Schedules -> Mo,Di,Mi,Do,Fr - Time
