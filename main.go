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

	http.HandleFunc("/off-net", disableNetAdpt)

	err := http.ListenAndServe(":5222", nil)
	checkError(err)
}

func disableNetAdpt(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)
	interval, err := strconv.Atoi(params["interval"])
	checkError(err)
	switchNetAdpt(interval, params["adptName"])
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

func checkError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}