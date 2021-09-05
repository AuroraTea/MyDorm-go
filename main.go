package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/test", postTest)
	http.HandleFunc("/off-net", disableNetAdpt)

	err := http.ListenAndServe(":5222", nil)
	checkError(err)
}

func postTest(w http.ResponseWriter, r *http.Request) {
	httpCORS(w, "*")
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var params map[string]string
		decoder.Decode(&params)

		interval, err := strconv.Atoi(params["interval"])
		checkError(err)
		fmt.Println(params["adptName"])
		fmt.Println(interval)
	}
}

func disableNetAdpt(w http.ResponseWriter, r *http.Request) {
	httpCORS(w, "*")
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var params map[string]string
		decoder.Decode(&params)
		interval, err := strconv.Atoi(params["interval"])
		checkError(err)

		switchNetAdpt(interval, params["adptName"])
	}
}

func switchNetAdpt(interval int, adptName string) {
	err :=exec.Command("netsh", "interface", "set", "interface", adptName, "disabled").Run()
	checkError(err)
	time.Sleep(time.Duration(interval) * time.Second)
	err = exec.Command("netsh", "interface", "set", "interface", adptName, "enabled").Run()
	checkError(err)
}

func getNetAdpt()  {
	out, err := exec.Command("netsh", "interface", "show", "interface").Output()
	checkError(err)
	fmt.Println(string(out))
}

func httpCORS(w http.ResponseWriter, url string) {
	w.Header().Set("Access-Control-Allow-Origin", url)
	w.Header().Add("Access-Control-Allow-Headers", "Access-Token, Content-Type")
	w.Header().Set("content-type", "application/json")
}

func httpMethod(r *http.Request, method string)  {
	if r.Method != method {
		return
	}
}


func checkError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}