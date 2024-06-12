package main

import (
	"io"
	"log"
	"net/http"

	"github.com/easysoft/gitfox-webhooks/gitfox"

	"github.com/davecgh/go-spew/spew"
)

// 接受所有webhook请求,并打印请求内容
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	log.Printf("request body: %s", body)
	spew.Dump(r.Header)
	hook, _ := gitfox.New()
	payload, err := hook.Parse(r, gitfox.BranchUpdatedEvent, gitfox.BranchCreatedEvent, gitfox.BranchDeletedEvent, gitfox.TagCreatedEvent)
	if err != nil {
		if err == gitfox.ErrEventNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Printf("parse webhook failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	switch payload := payload.(type) {
	case gitfox.BranchUpdatedPayload:
		log.Printf("branch updated payload: %v", payload)
	case gitfox.BranchCreatedPayload:
		log.Printf("branch created payload: %v", payload)
	case gitfox.BranchDeletedPayload:
		log.Printf("branch deleted payload: %v", payload)
	case gitfox.TagCreatedPayload:
		log.Printf("tag created payload: %v", payload)
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status": "ok"}`))
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/webhook", webhookHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
