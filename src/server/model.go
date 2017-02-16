package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

var (
	driverName     string
	dataSourceName string
)

type AQIData struct {
	Id   int
	City string
	Aqi  int
	Time string
}

type AQIDatas []AQIData

func init() {
	driverName = "mysql"
	dataSourceName = "root:123456@tcp(127.0.0.1:3306)/testcity?charset=utf8"
}

func Conn() *sql.DB {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err.Error())
		fmt.Println(err.Error())
	}
	return db
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func OneCitySingleData(t string, location string) (data AQIData, msg modelErr) {
	db := Conn()
	defer db.Close()
	msg.Exist = false

	// "\"防止引号转义
	rows, err := db.Query("SELECT * FROM testcity.aqi WHERE id=" + t + " AND city=\"" + location + "\";")
	checkErr(err)
	defer rows.Close()

	var (
		id   int
		city string
		aqi  int
		time string
	)

	for rows.Next() {
		err := rows.Scan(&id, &city, &aqi, &time)
		checkErr(err)
		data = AQIData{id, city, aqi, time}
		//datas = append(datas, data)
	}

	if data.City == "" || err != nil {
		msg.Exist = true
	}
	return
}

func OneCityAllDataToday(t string, location string) (datas AQIDatas, err error) {
	db := Conn()
	defer db.Close()

	tmp := time.Now().Hour()
	for i := 0; i <= tmp; i++ {
		h := t + strconv.Itoa(i)

		rows, err := db.Query("SELECT * FROM testcity.aqi WHERE id=" + h + " AND city=\"" + location + "\";")
		checkErr(err)
		defer rows.Close()

		var (
			id   int
			city string
			aqi  int
			time string
		)

		for rows.Next() {
			err := rows.Scan(&id, &city, &aqi, &time)
			checkErr(err)
			data := AQIData{id, city, aqi, time}
			datas = append(datas, data)
		}
	}
	return
}
