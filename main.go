package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

func main() {
	http.HandleFunc("/off-net", disableNetAdpt)

	err := http.ListenAndServe(":5222", nil)
	checkError(err)
	//showNet()
}

func disableNetAdpt(w http.ResponseWriter, r *http.Request) {
	_, err := exec.Command("netsh", "interface", "set", "interface", "以太网", "disabled").Output()
	checkError(err)
	time.Sleep(6 * time.Second)
	_, err = exec.Command("netsh", "interface", "set", "interface", "以太网", "enabled").Output()
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