// Copyright (c) 2024 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"

	"github.com/easysoft/gitfox-webhooks/gitfox"
)

// 接受所有webhook请求,并打印请求内容
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	hook, _ := gitfox.New()
	payload, err := hook.Parse(r,
		gitfox.BranchCreatedEvent, gitfox.BranchUpdatedEvent, gitfox.BranchDeletedEvent,
		gitfox.TagCreatedEvent, gitfox.TagDeletedEvent, gitfox.TagUpdatedEvent,
		gitfox.PullReqCreatedEvent, gitfox.PullReqReopenedEvent, gitfox.PullReqBranchUpdatedEvent,
		gitfox.PullReqClosedEvent, gitfox.PullReqCommentCreatedEvent, gitfox.PullReqMergedEvent,
		gitfox.PullReqReviewerCreatedEvent, gitfox.PullReqReviewerDeletedEvent, gitfox.PullReqReviewSubmittedEvent,
		gitfox.PullReqRequiredChecksPassedEvent,
	)
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
	case gitfox.BranchCreatedPayload:
		log.Printf("branch created: %v", payload)
	case gitfox.BranchUpdatedPayload:
		log.Printf("branch updated: %v", payload)
	case gitfox.BranchDeletedPayload:
		log.Printf("branch deleted: %v", payload)
	case gitfox.TagCreatedPayload:
		log.Printf("tag created: %v", payload)
	case gitfox.TagDeletedPayload:
		log.Printf("tag deleted: %v", payload)
	case gitfox.TagUpdatedPayload:
		log.Printf("tag updated: %v", payload)
	case gitfox.PullReqCreatedPayload:
		log.Printf("pull request created: %v", payload)
	case gitfox.PullReqReopenedPayload:
		log.Printf("pull request reopened: %v", payload)
	case gitfox.PullReqBranchUpdatedPayload:
		log.Printf("pull request branch updated: %v", payload)
	case gitfox.PullReqClosedPayload:
		log.Printf("pull request closed: %v", payload)
	case gitfox.PullReqCommentPayload:
		log.Printf("pull request comment created: %v", payload)
	case gitfox.PullReqMergedPayload:
		log.Printf("pull request merged: %v", payload)
	case gitfox.PullReqReviewerCreatedPayload:
		log.Printf("pull request reviewer created: %v", payload)
	case gitfox.PullReqReviewerDeletedPayload:
		log.Printf("pull request reviewer deleted: %v", payload)
	case gitfox.PullReqReviewSubmittedPayload:
		log.Printf("pull request review submitted: %v", payload)
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status": "ok"}`))
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/webhook", webhookHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
