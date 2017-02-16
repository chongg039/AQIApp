package main

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron"
	"regexp"
	//"strconv"
	"strings"
	"time"
)

var (
	driverName     string
	dataSourceName string
)

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

func Store() {
	n := time.Now()
	db := Conn()
	defer db.Close()
	// 调用事务插入数据，避免单条插入
	tx, err := db.Begin()
	if err != nil {
		panic(err.Error())
		fmt.Println(err.Error())
	}

	stmt, err := tx.Prepare("INSERT aqi SET id=?, city=?, aqi=?, time=?")
	if err != nil {
		panic(err.Error())
		fmt.Println(err.Error())
	}

	doc, err := goquery.NewDocument("http://www.pm25.com/rank.html")
	if err != nil {
		panic(err.Error())
		fmt.Println(err.Error())
	}

	timeTemp := doc.Find(".rank_banner_right span").Text()
	re, _ := regexp.Compile("[0-9]")
	times := re.FindAllString(timeTemp, -1)
	timeNow := strings.Join(times[:10], "")

	doc.Find(".rank_box").Find(".pj_area_data_item").Each(func(i int, s *goquery.Selection) {
		location := s.Find(".pjadt_location").Text()
		aqi := s.Find(".pjadt_aqi").Text()

		i++
		_, err := stmt.Exec(timeNow, location, aqi, timeTemp)
		if err != nil {
			panic(err.Error())
			fmt.Println(err.Error())
		}
	})
	// 出异常回滚
	defer tx.Rollback()

	// 提交事务
	tx.Commit()

	// 计时器模块，需要注释掉
	m := time.Since(n).Seconds()
	fmt.Println("[OK]", "Use", m, "'s: Store AQI datas at", time.Now())
}

func Exist() {
	db := Conn()
	defer db.Close()
	//t := time.Now().Format("2006010215")
	doc, err := goquery.NewDocument("http://www.pm25.com/rank.html")
	if err != nil {
		panic(err.Error())
		fmt.Println(err.Error())
	}
	timeTemp := doc.Find(".rank_banner_right span").Text()
	re, _ := regexp.Compile("[0-9]")
	times := re.FindAllString(timeTemp, -1)
	t := strings.Join(times[:10], "")
	// id, err := strconv.Atoi(t)
	// if err != nil {
	// 	panic(err.Error())
	// 	fmt.Println(err.Error())
	// }
	// idStr := strconv.Itoa(id)
	rows, err := db.Query("SELECT * FROM testcity.aqi WHERE id=" + t + " LIMIT 1")
	if err != nil {
		panic(err.Error())
		fmt.Println(err.Error())
	}
	//fmt.Println(rows)
	//sign := false
	var idtime, aqi int
	var city, time string

	for rows.Next() {
		// var id, aqi int
		// var city, time string
		rows.Columns()
		err = rows.Scan(&idtime, &city, &aqi, &time)
		if err != nil {
			panic(err.Error())
			fmt.Println(err.Error())
		}
	}
	// 选择检测aqi而不是id，因为清空数据后仍可能留有id
	if aqi > 0 {
		// sign = true
		//fmt.Println("Already have this time data,can't store")
		return
	}
	Store()
	//return sign
}

// func Ticker(id int) {
// 	sign := false
// 	for sign == false {
// 		Store()
// 		ex := Exist(id)
// 		sign = ex
// 	}
// }

func main() {
	c := cron.New()
	// t := time.Now().Format("2006010215")
	c.AddFunc("@every 1s", func() {
		if time.Now().Minute() < 1 {
			//fmt.Println(time.Now())
			Exist()
		}
		//return
	})
	c.Start()
	select {}
}
