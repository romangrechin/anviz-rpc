package main

import (
	"fmt"
	"github.com/romangrechin/anviz-rpc/anviz/device"
	"github.com/romangrechin/anviz-rpc/anviz/models"
	"github.com/romangrechin/anviz-rpc/api"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ACK_SUCCESS = 0x00
	ACK_FAIL    = 0x01
)

func main() {
	api.RunServer(":8081")
	os.Exit(0)

	if len(os.Args) < 3 {
		fmt.Println("Usage: anviz-info.exe host[:port] scanner-id (Default port: 5010)")
		return
	}

	port := 5010
	host := os.Args[1]

	if !strings.Contains(host, ":") {
		host = fmt.Sprintf("%s:%d", host, port)
	}

	deviceId, err := strconv.Atoi(os.Args[2])
	if err != nil || deviceId < 0 || deviceId > 999999999 {
		log.Println("Invalid device id")
		return
	}

	dev := device.New(host, 5*time.Second, 1*time.Second)
	err = dev.Connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println(dev.Id())
	log.Println(dev.GetInfo())
	log.Println(dev.GetFactoryInfoCode())

	first := true
	count := 0
	/*
		for {
			recs, err := dev.GetTaRecords(first, false, false)
			if err != nil{
				break
			}

			count = recs.Length()
			for _, rec := range recs.Get(){
				log.Println(rec.DateTime, rec.RecordType)
			}

			first = false

			if count < models.UsersMaxResults{
				break
			}
		}
	*/

	first = true
	count = 0
	max := models.AnsiUsersMaxResults
	if dev.IsUnicode() {
		max = models.UnicodeUsersMaxResults
	}

	for {
		log.Println("Start: ", time.Now())
		recs, err := dev.GetUsers(first, false)
		if err != nil {
			log.Println(err)
			break
		}

		count = recs.Length()
		for _, rec := range recs.Get() {
			log.Println(rec.UserCode, rec.CardCode, rec.Name)
		}

		first = false
		log.Println("Stop: ", time.Now())
		if count < max {
			break
		}
	}

	os.Exit(0)
}
