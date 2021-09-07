package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"reflect"
	"strconv"
	"time"
)

func main() {
	go checkIP()

	http.HandleFunc("/test", postTest)
	http.HandleFunc("/off-net", disableNetAdpt)

	http.HandleFunc("/get-ip", updateIP)

	err := http.ListenAndServe(":5222", nil)
	checkError(err)
}

func updateIP(w http.ResponseWriter, r *http.Request) {
	httpCORS(w, "*")
	if r.Method == "GET" {
		newIP,err := getLocalIPv4s()
		checkError(err)

		newIPString := ""
		for _,str := range newIP {
			newIPString += str+"\n"
		}

		fmt.Fprintf(w, newIPString)
	}
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

func checkIP()  {
	t := time.NewTicker(2 * time.Second)
	defer t.Stop()
	currentIP,err := getLocalIPv4s()
	checkError(err)
	for {
		<- t.C
		fmt.Println("Ticker running...")
		newIP,err := getLocalIPv4s()
		checkError(err)
		if !reflect.DeepEqual(newIP, currentIP){
			currentIP = newIP
		}
	}
}

func getLocalIPv4s() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	return ips, nil
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