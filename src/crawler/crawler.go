package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron"
	"gopkg.in/mgo.v2"
	"log"
	"regexp"
	"strings"
	"time"
)

const (
	URL = "127.0.0.1:27017"
)

type AirDataItem struct {
	City string `json: city`
	AQI  string `json: aqi`
	Time string `json: time`
}

type AirDatas struct {
	DataItems []AirDataItem
}

func GetAndStore() {
	session, err := mgo.Dial(URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	var airDataItem AirDataItem
	doc, err := goquery.NewDocument("http://www.pm25.com/rank.html")

	if err != nil {
		log.Fatal(err)
	}

	// use time as collection
	timeTemp := doc.Find(".rank_banner_right span").Text()
	re, _ := regexp.Compile("[0-9]")
	times := re.FindAllString(timeTemp, -1)
	timeNow := strings.Join(times[:10], "")

	db := session.DB("airDatas")
	collection := db.C(timeNow)

	doc.Find(".pj_area_data_item").Each(func(i int, s *goquery.Selection) {
		location := s.Find(".pjadt_location").Text()
		aqi := s.Find(".pjadt_aqi").Text()

		airDataItem.City = location
		airDataItem.AQI = aqi
		airDataItem.Time = timeTemp
		err = collection.Insert(airDataItem)
		if err != nil {
			log.Fatal(err)
		}
	})
}

func AQITicker() {
	c := cron.New()
	t := time.Now()
	c.AddFunc("@every 1m", func() {
		log.Println("Get and store", t.Local(), "data")
		GetAndStore()
	})
	c.Start()
	select {}
}

func main() {
	log.Println("Start crawler engine...")
	AQITicker()
}
