// Copyright (c) 2024 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package gitfox

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// parse errors
var (
	ErrEventNotSpecifiedToParse             = errors.New("no Event specified to parse")
	ErrInvalidHTTPMethod                    = errors.New("invalid HTTP Method")
	ErrMissingGitFoxEventHeader             = errors.New("missing X-Gitfox-Event Header")
	ErrMissingGitFoxTriggerHeader           = errors.New("missing X-Gitfox-Trigger Header")
	ErrMissingGitFoxWebhookParentTypeHeader = errors.New("missing Webhook-Parent-Type Header")
	ErrMissingGitFoxSignatureHeader         = errors.New("missing X-Gitfox-Signature Header")
	ErrEventNotFound                        = errors.New("event not defined to be parsed")
	ErrParsingPayload                       = errors.New("error parsing payload")
	ErrHMACVerificationFailed               = errors.New("HMAC verification failed")
)

type HookEventType string

const (
	BranchCreatedEvent HookEventType = "branch_created"
	BranchUpdatedEvent HookEventType = "branch_updated"
	BranchDeletedEvent HookEventType = "branch_deleted"

	TagCreatedEvent HookEventType = "tag_created"
	TagDeletedEvent HookEventType = "tag_deleted"
	TagUpdatedEvent HookEventType = "tag_updated"

	PullReqCreatedEvent        HookEventType = "pullreq_created"
	PullReqReopenedEvent       HookEventType = "pullreq_reopened"
	PullReqBranchUpdatedEvent  HookEventType = "pullreq_branch_updated"
	PullReqClosedEvent         HookEventType = "pullreq_closed"
	PullReqCommentCreatedEvent HookEventType = "pullreq_comment_created"
	PullReqCommentUpdatedEvent HookEventType = "pullreq_comment_updated"
	// PullReqCommentStatusUpdated HookEventType = "pullreq_comment_status_updated"
	PullReqMergedEvent               HookEventType = "pullreq_merged"
	PullReqReviewerCreatedEvent      HookEventType = "pullreq_reviewer_created"
	PullReqRequiredChecksPassedEvent HookEventType = "pullreq_required_checks_passed"
	PullReqReviewerDeletedEvent      HookEventType = "pullreq_reviewer_deleted"
	PullReqReviewSubmittedEvent      HookEventType = "pullreq_review_submitted"
	PullReqUpdatedEvent              HookEventType = "pullreq_updated"
)

// Option is a configuration option for the webhook
type Option func(*Webhook) error

// Options is a namespace var for configuration options
var Options = WebhookOptions{}

// WebhookOptions is a namespace for configuration option methods
type WebhookOptions struct{}

// Secret registers the GitLab secret
func (WebhookOptions) Secret(secret string) Option {
	return func(hook *Webhook) error {
		hook.secret = secret
		return nil
	}
}

// Webhook instance contains all methods needed to process events
type Webhook struct {
	secret string
}

// New creates and returns a WebHook instance denoted by the Provider type
func New(options ...Option) (*Webhook, error) {
	hook := new(Webhook)
	for _, opt := range options {
		if err := opt(hook); err != nil {
			return nil, errors.New("Error applying Option")
		}
	}
	return hook, nil
}

// Parse verifies and parses the events specified and returns the payload object or an error
func (hook Webhook) Parse(r *http.Request, events ...HookEventType) (interface{}, error) {
	defer func() {
		_, _ = io.Copy(io.Discard, r.Body)
		_ = r.Body.Close()
	}()

	if len(events) == 0 {
		return nil, ErrEventNotSpecifiedToParse
	}
	if r.Method != http.MethodPost {
		return nil, ErrInvalidHTTPMethod
	}
	event := r.Header.Get("X-Gitfox-Trigger")
	if len(event) == 0 {
		return nil, ErrMissingGitFoxTriggerHeader
	}
	gitfoxEvent := HookEventType(event)
	var found bool
	for _, evt := range events {
		if evt == gitfoxEvent {
			found = true
			break
		}
	}
	// event not defined to be parsed
	if !found {
		return nil, ErrEventNotFound
	}
	payload, err := io.ReadAll(r.Body)
	if err != nil || len(payload) == 0 {
		return nil, ErrParsingPayload
	}
	// If we have a Secret set, we should check the MAC
	if len(hook.secret) > 0 {
		signature := r.Header.Get("X-Gitfox-Signature")
		if len(signature) == 0 {
			return nil, ErrMissingGitFoxSignatureHeader
		}
		sig256 := hmac.New(sha256.New, []byte(hook.secret))
		_, _ = io.Writer(sig256).Write([]byte(payload))
		expectedMAC := hex.EncodeToString(sig256.Sum(nil))

		if !hmac.Equal([]byte(signature), []byte(expectedMAC)) {
			return nil, ErrHMACVerificationFailed
		}
	}
	switch gitfoxEvent {
	case BranchCreatedEvent:
		var pl BranchCreatedPayload
		err = json.Unmarshal([]byte(payload), &pl)
		return pl, err
	case BranchUpdatedEvent:
		var pl BranchUpdatedPayload
		err = json.Unmarshal([]byte(payload), &pl)
		return pl, err
	case BranchDeletedEvent:
		var pl BranchDeletedPayload
		err = json.Unmarshal([]byte(payload), &pl)
		return pl, err
	case TagCreatedEvent:
		var pl TagCreatedPayload
		err = json.Unmarshal([]byte(payload), &pl)
		return pl, err
	case TagUpdatedEvent:
		var pl TagUpdatedPayload
		err = json.Unmarshal([]byte(payload), &pl)
		return pl, err
	case TagDeletedEvent:
		var pl TagDeletedPayload
		err = json.Unmarshal([]byte(payload), &pl)
		return pl, err
	case PullReqCreatedEvent:
		var pl PullReqCreatedPayload
		err = json.Unmarshal([]byte(payload), &pl)
		return pl, err
	case PullReqReopenedEvent:
		var pl PullReqReopenedPayload
		err = json.Unmarshal([]byte(payload), &pl)
		return pl, err
	case PullReqBranchUpdatedEvent:
		var pl PullReqBranchUpdatedPayload
		err = json.Unmarshal([]byte(payload), &pl)
		return pl, err
	case PullReqClosedEvent:
		var pl PullReqClosedPayload
		err = json.Unmarshal([]byte(payload), &pl)
		return pl, err
	case PullReqCommentCreatedEvent:
		var pl PullReqCommentPayload
		err = json.Unmarshal([]byte(payload), &pl)
		return pl, err
	case PullReqCommentUpdatedEvent:
		var pl PullReqCommentUpdatedPayload
		err = json.Unmarshal([]byte(payload), &pl)
		return pl, err
	case PullReqMergedEvent:
		var pl PullReqMergedPayload
		err = json.Unmarshal([]byte(payload), &pl)
		return pl, err
	case PullReqUpdatedEvent:
		var pl PullReqUpdatedPayload
		err = json.Unmarshal([]byte(payload), &pl)
		return pl, err
	case PullReqReviewerCreatedEvent, PullReqRequiredChecksPassedEvent:
		var pl PullReqReviewerCreatedPayload
		err = json.Unmarshal([]byte(payload), &pl)
		return pl, err
	case PullReqReviewerDeletedEvent:
		var pl PullReqReviewerDeletedPayload
		err = json.Unmarshal([]byte(payload), &pl)
		return pl, err
	case PullReqReviewSubmittedEvent:
		var pl PullReqReviewSubmittedPayload
		err = json.Unmarshal([]byte(payload), &pl)
		return pl, err
	default:
		return nil, fmt.Errorf("event type %s not yet supported", gitfoxEvent)
	}
}
