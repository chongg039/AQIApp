package main

type AirDataItem struct {
	City string `json: city`
	AQI  string `json: aqi`
	Time string `json: time`
}

type AirDatas struct {
	DataItems []AirDataItem
}
