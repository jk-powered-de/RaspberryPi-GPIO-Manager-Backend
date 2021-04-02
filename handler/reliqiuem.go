package handler

import (
	"github.com/stianeikeland/go-rpio"
	"log"
	"raspberrypi-gpio-manager-backend/db"
	"raspberrypi-gpio-manager-backend/model"
	"time"
)

func Core() {

	prepareJobs()

	controlGpioPins()
}

func controlGpioPins() {
	err := rpio.Open()
	if err != nil {
		log.Printf("unable to use gpiopins: %s", err.Error())
		return
	}
	defer rpio.Close()

	for _, namedGpioPin := range getPinsNeedToBeStarted() {
		pin := rpio.Pin(namedGpioPin.Pin)

		pin.Output()
		pin.High()
	}

	for _, namedGpioPin := range getPinsNeedToBeStopped() {
		pin := rpio.Pin(namedGpioPin.Pin)

		pin.Input()
		// pin.Output()
		// pin.Low()
	}
}

func prepareJobs() {
	for _, job := range getJobsNeedToBeStop() {
		turnOffNamedGpioPins(getNamedGpioPin(job.NamedPinId).ID)
		changeState(job, "done")
	}

	for _, job := range getJobsNeedToBeStart() {
		turnOnNamedGpioPins(getNamedGpioPin(job.NamedPinId).ID)
		changeState(job, "running")
	}
}

func turnOnNamedGpioPins(id uint) {
	var namedGpioPin model.NamedPin

	raw := db.Connection.Debug().Where("id = ?", id).Model(&namedGpioPin).Update("state", "on")
	if raw.Error != nil {
		printError(raw.Error, 503, "Database Service unavailable", nil)
		return
	}

}

func turnOffNamedGpioPins(id uint) {
	var namedGpioPin model.NamedPin

	raw := db.Connection.Debug().Where("id = ?", id).Model(&namedGpioPin).Update("state", "off")
	if raw.Error != nil {
		printError(raw.Error, 503, "Database Service unavailable", nil)
		return
	}

}

func changeState(job model.Job, state string) {

	job.State = state

	raw := db.Connection.Debug().Where("id = ?", job.ID).Model(&job).Update("state", state)
	if raw.Error != nil {
		log.Fatalln(raw.Error)
		return
	}
}

func StartInterval() {
	for {
		time.Sleep(10 * time.Second)
		go Core()
	}
}

func getPinsNeedToBeStarted() []model.NamedPin {

	var namedGpioPins []model.NamedPin
	raw := db.Connection.Debug().Where("state = ?", "on").Find(&namedGpioPins)
	if raw.Error != nil {
		log.Fatalln(raw.Error)
	}

	return namedGpioPins
}

func getPinsNeedToBeStopped() []model.NamedPin {

	var namedGpioPins []model.NamedPin
	raw := db.Connection.Debug().Where("state = ?", "off").Find(&namedGpioPins)
	if raw.Error != nil {
		log.Fatalln(raw.Error)
	}

	return namedGpioPins
}

func getJobsNeedToBeStart() []model.Job {

	var jobs []model.Job
	raw := db.Connection.Debug().Where("start_time < ?", time.Now().Unix()).Where("state != 'running'").Where("state != 'done'").Where("end_time > ?", time.Now().Unix()).Find(&jobs)
	if raw.Error != nil {
		log.Fatalln(raw.Error)
	}

	return jobs
}

func getJobsNeedToBeStop() []model.Job {

	var jobs []model.Job
	raw := db.Connection.Debug().Where("state = 'running'").Where("end_time < ?", time.Now().Unix()).Find(&jobs)
	if raw.Error != nil {
		log.Fatalln(raw.Error)
	}

	return jobs
}

func getNamedGpioPin(namedGpioPinId int) model.NamedPin {

	var namedGpioPin model.NamedPin
	raw := db.Connection.Debug().Where("id = ?", namedGpioPinId).First(&namedGpioPin)
	if raw.Error != nil {
		log.Fatalln(raw.Error)
	}

	return namedGpioPin
}
