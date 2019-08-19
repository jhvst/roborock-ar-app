package main

import (
	"log"
	"net/http"
	"os/exec"
)

func handler(w http.ResponseWriter, r *http.Request) {
	exec.Command("/bin/bash", "./switch.sh").Run()
	return
}

func handler2(w http.ResponseWriter, r *http.Request) {
	exec.Command("/bin/bash", "./on.sh").Run()
	return
}

func media(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", `*`)
	data, _ := exec.Command("/bin/bash", "./media.sh").Output()
	w.Write(data)
}

func pause(w http.ResponseWriter, r *http.Request) {
	exec.Command("/bin/bash", "./pause.sh").Run()
	return
}

func main() {
	http.HandleFunc("/off", handler)
	http.HandleFunc("/on", handler2)
	http.HandleFunc("/media", media)
	http.HandleFunc("/pause", pause)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
