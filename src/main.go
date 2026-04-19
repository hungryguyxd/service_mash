package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var path = "/app/logs/app.log"

type Message struct {
	Message string `json:"message"`
}

func getMessage() string {
	data, err := os.ReadFile("/app/config/message")
	if err != nil {
		return "попался разбойник"
	}
	return string(data)
}

func getPort() string {
	data, err := os.ReadFile("/app/config/port")
	if err != nil {
		return "ой ой ой"
	}
	return string(data)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	podName := os.Getenv("HOSTNAME")
	fmt.Fprintf(w, "pod: %s, message: %s\n", podName, getMessage())
	fmt.Fprintln(w, getMessage())
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	var mes Message
	err := json.NewDecoder(r.Body).Decode(&mes)
	if err != nil {
		http.Error(w, "давай по новой", http.StatusBadRequest)
		return
	}
	log.Println(mes.Message)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "печалька", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	f.WriteString(mes.Message + "\n")
	w.WriteHeader(http.StatusOK)
}

func logsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "404 грустно", http.StatusNotFound)
		return
	}
	w.Write(data)
}

func main() {
	port := getPort()
	os.MkdirAll("/app/logs", os.ModePerm)
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/log", logHandler)
	http.HandleFunc("/logs", logsHandler)
	log.Printf("cтратуем сервер на с[порт]ике):%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
