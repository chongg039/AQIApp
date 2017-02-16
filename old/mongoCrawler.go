package main

import (
	//"fmt"
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
	//sign = true
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

	n := time.Now()
	doc.Find(".rank_box").Find(".pj_area_data_item").Each(func(i int, s *goquery.Selection) {
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
	m := time.Since(n).Nanoseconds()
	log.Println(m)
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
	// AQITicker()
	GetAndStore()
}

// func main() {

// //time.NewTicker生成一个ticket，它包含一个管道channel C，然后每个相应的时间间隔，会向管道发送数据。
// //我们使用for range遍历管道，就实现了间隔时间定时执行的问题。
// ticker := time.NewTicker(time.Second * 1)
// go func() {
// 	for t := range ticker.C {
// 		// fmt.Println("Tick at", t)

// 		if time.Now().Minute() <= 7 {
// 			re := GetAndStore()
// 			if re == false {
// 				continue
// 			}
// 			fmt.Println("Store data at ", t)
// 		}
// 		continue
// 	}
// }()

// //和timers一样，tickes也可以被停止，停止后，管道就不会接受值了。
// time.Sleep(time.Hour * 87600)
// ticker.Stop()
// fmt.Println("Ticker stopped")
// }
