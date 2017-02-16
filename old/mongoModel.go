package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	URL = "127.0.0.1:27017"
)

func OneCitySingleData(city string, time string) (AirDataItem, modelErr) {
	session, err := mgo.Dial(URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	db := session.DB("airDatas")
	collection := db.C(time)

	result := AirDataItem{}
	errMessage := modelErr{}
	err = collection.Find(bson.M{"city": city}).One(&result)
	if err != nil || result.City == "" {
		// log.Fatal(err)
		errMessage.Exist = true
	} else {
		errMessage.Exist = false
	}
	// errMessage.Exist = false
	return result, errMessage
}

func OneCityAllDataToday(city string, day string) AirDatas {
	session, err := mgo.Dial(URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	db := session.DB("airDatas")

	// re, _ := regexp.Compile("[0-9]")
	// timeTemp := re.FindAllString(time, -1)
	// hour := strings.Join(timeTemp[8:10], "")
	hour := strings.Replace(time.Now().String()[11:13], "-", "", -1)
	// timeToday := strings.Join(timeTemp[:8], "")

	var timeHistory string
	var result AirDatas
	// var hourToInt int

	hourToInt, err := strconv.Atoi(hour)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i <= hourToInt; i++ {
		if i < 10 {
			timeHistory = day + "0" + strconv.Itoa(i)
		} else {
			timeHistory = day + strconv.Itoa(i)
		}
		collection := db.C(timeHistory)

		temp := AirDataItem{}
		err = collection.Find(bson.M{"city": city}).One(&temp)
		if err != nil {
			// result.Reminder = append(result.Reminder, "暂时没有这个时段的数据！")
			// result.DataItems = append(result.DataItems, nil)
			continue
		}
		result.DataItems = append(result.DataItems, temp)
	}
	return result
}
