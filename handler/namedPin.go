package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stianeikeland/go-rpio"
	"log"
	"raspberrypi-gpio-manager-backend/db"
	"raspberrypi-gpio-manager-backend/model"
)

func FindAllNamedPins(c *gin.Context) {

	var namedPins []model.NamedPin
	raw := db.Connection.Debug().Find(&namedPins)
	if raw.Error != nil {
		printError(raw.Error, 503, "Database Service unavailable", c)
		return
	}

	c.JSON(200, gin.H{
		"code":    "200",
		"message": "",
		"data":    &namedPins,
	})
}

func CreateNamedPin(c *gin.Context) {
	var namedPins []model.NamedPin
	var createNamedPin model.NamedPin
	if err := c.ShouldBindJSON(&createNamedPin); err != nil {
		printError(err, 400, "Malformed Data", c)
		return
	}
	createNamedPin.State = "off"

	raw := db.Connection.Debug().Where("pin = ? ", createNamedPin.Pin).First(&namedPins)

	if raw.RowsAffected != 0 {
		printError(raw.Error, 400, "Pin is already set", c)
		return
	}

	raw = db.Connection.Debug().Where("name = ? ", createNamedPin.Name).First(&namedPins)

	if raw.RowsAffected != 0 {
		printError(raw.Error, 400, "Name already exist", c)
		return
	}

	raw = db.Connection.Debug().Create(&createNamedPin)
	if raw.Error != nil {
		fmt.Println(raw.Error)
		printError(raw.Error, 503, "Database Service unavailable", c)
		return
	}

	c.JSON(201, gin.H{
		"code":    "201",
		"message": "Created",
	})
}

func TurnOnNamedPinByID(c *gin.Context) {
	id := c.Param("id")

	var namedPin model.NamedPin

	raw := db.Connection.Debug().Where("id = ?", id).Model(&namedPin).Update("state", "on").First(&namedPin)
	if raw.Error != nil {
		printError(raw.Error, 503, "Database Service unavailable", c)
		return
	}

	err := rpio.Open()
	if err != nil {
		log.Printf("unable to use gpiopins: %s", err.Error())
		return
	}
	defer rpio.Close()

	pin := rpio.Pin(namedPin.Pin)
	pin.Output()
	pin.High()

	c.JSON(201, gin.H{
		"code":    "201",
		"message": "Turned pin on",
	})
}

func TurnOffNamedPinByID(c *gin.Context) {
	id := c.Param("id")

	var namedPin model.NamedPin

	raw := db.Connection.Debug().Where("id = ?", id).Model(&namedPin).Update("state", "off").First(&namedPin)
	if raw.Error != nil {
		printError(raw.Error, 503, "Database Service unavailable", c)
		return
	}

	err := rpio.Open()
	if err != nil {
		log.Printf("unable to use gpiopins: %s", err.Error())
		return
	}
	defer rpio.Close()

	pin := rpio.Pin(namedPin.Pin)
	pin.Input()
	//pin.Output()
	//pin.Low()

	c.JSON(201, gin.H{
		"code":    "201",
		"message": "Turned pin off",
	})
}

func UpdateNamedPinByID(c *gin.Context) {
	id := c.Param("id")
	var (
		namedPin    model.NamedPin
		namedPins   []model.NamedPin
		newNamedPin model.NamedPinPatch
	)

	if err := c.ShouldBindJSON(&newNamedPin); err != nil {
		printError(err, 400, "Malformed Data", c)
		return
	}

	raw := db.Connection.Debug().Where("pin = ? ", newNamedPin.Pin).Not("id = ?", id).First(&namedPins)

	if raw.RowsAffected != 0 {
		printError(raw.Error, 400, "Pin is already set", c)
		return
	}

	raw = db.Connection.Debug().Where("name = ? ", newNamedPin.Name).Not("id = ?", id).First(&namedPins)

	if raw.RowsAffected != 0 {
		printError(raw.Error, 400, "Name already exist", c)
		return
	}

	raw = db.Connection.Debug().Where("id = ?", id).First(&namedPin)
	if raw.Error != nil {
		printError(raw.Error, 503, "Database Service unavailable", c)
		return
	}

	log.Print(&namedPins)
	log.Print(&newNamedPin)
	raw = db.Connection.Debug().Where("id = ?", id).Model(&namedPin).Updates(&newNamedPin)
	if raw.Error != nil {
		printError(raw.Error, 503, "Database Service unavailable", c)
		return
	}

	if namedPin.State == "on" {
		err := rpio.Open()
		if err != nil {
			log.Printf("unable to use gpiopins: %s", err.Error())
			return
		}
		defer rpio.Close()
		pin := rpio.Pin(namedPin.Pin)
		pin.Input()
		//pin.Output()
		//pin.Low()

		pin = rpio.Pin(newNamedPin.Pin)
		pin.Output()
		pin.High()
	}

	if raw.RowsAffected == 0 {
		c.JSON(200, gin.H{
			"code":    "200",
			"message": "Nothing to update",
		})
		return
	}

	c.JSON(201, gin.H{
		"code":    "201",
		"message": "Updated",
	})
}

func DeleteNamedPinByID(c *gin.Context) {
	id := c.Param("id")

	var namedPin model.NamedPin
	raw := db.Connection.Debug().Where("id = ?", id).First(&namedPin)
	if raw.Error != nil {
		printError(raw.Error, 503, "Database Service unavailable", c)
		return
	}

	err := rpio.Open()
	if err != nil {
		log.Printf("unable to use gpiopins: %s", err.Error())
	}
	defer rpio.Close()

	pin := rpio.Pin(namedPin.Pin)
	pin.Input()
	//pin.Output()
	//pin.Low()

	var jobs []model.Job
	raw = db.Connection.Debug().Where("named_gpio_pin_id = ?", id).Find(&jobs)
	if raw.Error != nil {
		printError(raw.Error, 503, "Database Service unavailable", c)
		return
	}

	raw = db.Connection.Debug().Delete(&jobs, "named_gpio_pin_id = ?", id)
	if raw.Error != nil {
		printError(raw.Error, 503, "Database Service unavailable", c)
		return
	}

	raw = db.Connection.Debug().Delete(&namedPin)
	if raw.Error != nil {
		printError(err, 503, "Database Service unavailable", c)
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
