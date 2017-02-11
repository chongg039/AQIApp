package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func ReturnSperifiedData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	time := mux.Vars(r)["time"]
	city := mux.Vars(r)["city"]
	result, message := OneCitySingleData(city, time)
	// final, err := json.Marshal(result)
	// if err != nil {
	// 	fmt.Println("json err: ", err)
	// }
	// fmt.Fprintln(w, string(final))
	//length := len(string(result))
	if message.Exist == true {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			panic(err)
		}
		return
	}

	// Didn't find, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}

func ReturnNowData(w http.ResponseWriter, r *http.Request) {
	timeTemp := strings.Replace(time.Now().String()[0:13], "-", "", -1)
	timeNow := strings.Replace(timeTemp, " ", "", -1)
	city := mux.Vars(r)["city"]
	result, message := OneCitySingleData(city, timeNow)

	if message.Exist == true {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			panic(err)
		}
		return
	}

	// Didn't find, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}

func ReturnAllDataToday(w http.ResponseWriter, r *http.Request) {
	day := strings.Replace(time.Now().String()[0:10], "-", "", -1)
	city := mux.Vars(r)["city"]
	result := OneCityAllDataToday(city, day)
	final, err := json.Marshal(result)
	if err != nil {
		fmt.Println("json err: ", err)
	}
	fmt.Fprintln(w, string(final))
}
