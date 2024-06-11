package main

import (
	"io"
	"log"
	"net/http"
)

var (
	Version, BuildDate, GitBranch, GitCommitHash string
)

// 接受所有webhook请求,并打印请求内容
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("read request body failed: %v", err)
		return
	}
	log.Printf("request body: %s", body)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status": "ok"}`))
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/webhook", webhookHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
