package handler

import (
	"github.com/gin-gonic/gin"
	"raspberrypi-gpio-manager-backend/db"
	"raspberrypi-gpio-manager-backend/model"
	"strconv"
)

func FindAllJobs(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(-1)))
	if err != nil {
		printError(err, 500, "Cannot convert limit to integer", c)
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", strconv.Itoa(-1)))
	if err != nil {
		printError(err, 500, "Cannot convert offset to integer", c)
		return
	}

	var jobs []model.Job
	raw := db.Connection.Debug().Offset(offset).Limit(limit).Order("start_time asc").Find(&jobs)
	if raw.Error != nil {
		printError(raw.Error, 503, "Database Service unavailable", c)
		return
	}

	c.JSON(200, gin.H{
		"code":    "200",
		"message": "",
		"data":    &jobs,
	})
}

func FindAllJobsUndone(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(-1)))
	if err != nil {
		printError(err, 500, "Cannot convert limit to integer", c)
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", strconv.Itoa(-1)))
	if err != nil {
		printError(err, 500, "Cannot convert offset to integer", c)
		return
	}

	var jobs []model.Job
	raw := db.Connection.Debug().Offset(offset).Limit(limit).Order("start_time desc").Where("state != 'done'").Find(&jobs)
	if raw.Error != nil {
		printError(raw.Error, 503, "Database Service unavailable", c)
		return
	}

	c.JSON(200, gin.H{
		"code":    "200",
		"message": "",
		"data":    &jobs,
	})
}

func FindAllJobsDone(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(-1)))
	if err != nil {
		printError(err, 500, "Cannot convert limit to integer", c)
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", strconv.Itoa(-1)))
	if err != nil {
		printError(err, 500, "Cannot convert offset to integer", c)
		return
	}

	var jobs []model.Job
	raw := db.Connection.Debug().Offset(offset).Limit(limit).Order("start_time asc").Where("state = 'done'").Find(&jobs)
	if raw.Error != nil {
		printError(raw.Error, 503, "Database Service unavailable", c)
		return
	}

	c.JSON(200, gin.H{
		"code":    "200",
		"message": "",
		"data":    &jobs,
	})
}

func CreateJob(c *gin.Context) {
	var job model.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		printError(err, 400, "Malformed Data", c)
		return
	}

	raw := db.Connection.Debug().Create(&job)
	if raw.Error != nil {
		printError(raw.Error, 503, "Database Service unavailable", c)
		return
	}

	c.JSON(201, gin.H{
		"code":    "201",
		"message": "Created",
	})
}

func DeleteJobByNamedGpioPinId(c *gin.Context) {
	namedGpioPinId := c.Param("id")

	var jobs []model.Job
	raw := db.Connection.Debug().Where("named_gpio_pin_id = ?", namedGpioPinId).Find(&jobs)
	if raw.Error != nil {
		printError(raw.Error, 503, "Database Service unavailable", c)
		return
	}

	for _, job := range jobs {
		var namedGpioPin model.NamedPin

		raw = db.Connection.Debug().Where("id = ?", job.NamedPinId).Model(&namedGpioPin).Update("state", "off")
		if raw.Error != nil {
			printError(raw.Error, 503, "Database Service unavailable", nil)
			return
		}
	}

	raw = db.Connection.Debug().Delete(&jobs, "named_gpio_pin_id = ?", namedGpioPinId)
	if raw.Error != nil {
		printError(raw.Error, 503, "Database Service unavailable", c)
		return
	}
	if raw.RowsAffected == 0 {
		c.JSON(200, gin.H{
			"code":    "200",
			"message": "Nothing to delete",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    "200",
		"message": "Deleted",
	})
}

func DeleteJobByID(c *gin.Context) {
	id := c.Param("id")

	var job model.Job
	raw := db.Connection.Debug().Where("id = ?", id).First(&job)
	if raw.Error != nil {
		printError(raw.Error, 503, "Database Service unavailable", c)
		return
	}

	var namedGpioPin model.NamedPin

	raw = db.Connection.Debug().Where("id = ?", job.NamedPinId).Model(&namedGpioPin).Update("state", "off")
	if raw.Error != nil {
		printError(raw.Error, 503, "Database Service unavailable", nil)
		return
	}

	raw = db.Connection.Debug().Delete(&job, "id = ?", id)
	if raw.Error != nil {
		printError(raw.Error, 503, "Database Service unavailable", c)
		return
	}
	if raw.RowsAffected == 0 {
		c.JSON(200, gin.H{
			"code":    "200",
			"message": "Nothing to delete",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    "200",
		"message": "Deleted",
	})
}
