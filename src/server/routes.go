package main

import (
	"net/http"
	//"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"ReturnSperifiedData",
		"GET",
		"/aqi/{city:[\u4E00-\u9FA5]+}&{time:[0-9]+}",
		ReturnSperifiedData,
	},
	Route{
		"ReturnNowData",
		"GET",
		"/aqi/{city:[\u4E00-\u9FA5]+}&now",
		ReturnNowData,
	},
	Route{
		"ReturnAllDataToday",
		"GET",
		"/aqi/{city:[\u4E00-\u9FA5]+}&today",
		ReturnAllDataToday,
	},
}
