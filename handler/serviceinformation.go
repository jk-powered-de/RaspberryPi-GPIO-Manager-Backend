package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"raspberrypi-gpio-manager-backend/model"
)

func ServiceInfo(c *gin.Context) {

	configFile, err := os.Open("config/config.serviceinformation.json")
	if err != nil {
		log.Fatalln(err)
	}

	defer configFile.Close()
	decoder := json.NewDecoder(configFile)
	Service := model.ServiceInfo{}
	err = decoder.Decode(&Service)
	if err != nil {
		fmt.Println("error:", err)
	}

	message := Service.VersionState + " v." + Service.Version + " " + Service.Name + " - " + Service.VersionName

	c.JSON(200, gin.H{
		"code":          "200",
		"message":       message,
		"name":          Service.Name,
		"version":       Service.Version,
		"version_state": Service.VersionState,
		"version_name":  Service.VersionName,
		"author":        Service.Author,
		"contributors":  Service.Contributor,
	})
}
