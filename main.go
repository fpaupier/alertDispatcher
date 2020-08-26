package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
	"time"
)

const (
	DriverName = "sqlite3"
	DBPath     = "/Users/fpaupier/projects/pi-mask-detection/alert.db"
)

func main() {
	database, err := sql.Open(DriverName, DBPath)
	if err != nil {
		log.Fatalf("failed to open database: %s\n", err.Error())
	}
	var id int
	var createdAt string
	var deviceType string
	var deviceId string
	var deviceDeployedOn string
	var longitude float32
	var latitude float32
	var faceModelName string
	var faceModelThreshold float32
	var faceModelGuid string
	var maskModelName string
	var maskModelGuid string
	var maskModelThreshold float32
	var probability float32
	var imageFormat string
	var imageWidth int
	var imageHeight int
	var imageData []byte

	for {
		// read from a sqlite db
		rows, err := database.Query("SELECT id, created_at, device_type, device_id, device_deployed_on, longitude, latitude, face_model_name, face_model_guid, face_model_threshold, mask_model_name, mask_model_guid, mask_model_threshold, probability, image_format, image_width, image_height, image_data FROM alert WHERE sent = 0")
		if err != nil {
			log.Fatalf("failed to read rows from db: %v\n", err)
		}

		//	Process all records not sent
		for rows.Next() {
			err = rows.Scan(&id, &createdAt, &deviceType, &deviceId, &deviceDeployedOn, &longitude, &latitude, &faceModelName, &faceModelGuid, &faceModelThreshold, &maskModelName, &maskModelGuid, &maskModelThreshold, &probability, &imageFormat, &imageWidth, &imageHeight, &imageData)
			if err != nil {
				log.Fatalf("failed to scan row: %v\n", err)
			}
			fmt.Println(strconv.Itoa(id), ": ", createdAt, deviceType)
		}

		time.Sleep(1 * time.Second)
	}

	//	Create protobuf

	//	Send message to kafka queue

}
