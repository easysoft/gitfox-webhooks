// Copyright (c) 2024 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package gitfox

import "time"

type Repo struct {
	ID            int    `json:"id"`
	Path          string `json:"path"`
	Identifier    string `json:"identifier"`
	DefaultBranch string `json:"default_branch"`
	GitURL        string `json:"git_url"`
	UID           string `json:"uid,omitempty"`
}

type PrincipalInfo struct {
	ID          int    `json:"id"`
	UID         string `json:"uid"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Type        string `json:"type"`
	Created     int64  `json:"created"`
	Updated     int64  `json:"updated"`
}

type IdentityInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SignatureInfo struct {
	Identity IdentityInfo `json:"identity"`
	When     time.Time    `json:"when"`
}

type CommitInfo struct {
	Sha       string        `json:"sha"`
	Message   string        `json:"message"`
	Author    SignatureInfo `json:"author"`
	Committer SignatureInfo `json:"committer"`
	Added     []any         `json:"added"`
	Removed   []any         `json:"removed"`
	Modified  []string      `json:"modified"`
}

type ReferenceInfo struct {
	Name string `json:"name"`
	Repo Repo   `json:"repo"`
}

type BranchCreatedPayload ReferencePayload
type BranchUpdatedPayload ReferencePayload
type BranchDeletedPayload ReferencePayload

type TagCreatedPayload ReferencePayload
type TagUpdatedPayload ReferencePayload
type TagDeletedPayload ReferencePayload

type ReferencePayload struct {
	BaseSegment
	ReferenceSegment
	ReferenceDetailsSegment
	ReferenceUpdateSegment
}

type BaseSegment struct {
	Trigger   string        `json:"trigger"`
	Repo      Repo          `json:"repo"`
	Principal PrincipalInfo `json:"principal"`
}

type ReferenceSegment struct {
	Ref ReferenceInfo `json:"ref"`
}

type ReferenceDetailsSegment struct {
	SHA string `json:"sha"`

	HeadCommit *CommitInfo `json:"head_commit,omitempty"`

	Commits           *[]CommitInfo `json:"commits,omitempty"`
	TotalCommitsCount int           `json:"total_commits_count,omitempty"`

	// Deprecated
	Commit *CommitInfo `json:"commit,omitempty"`
}

type ReferenceUpdateSegment struct {
	OldSHA string `json:"old_sha"`
	Forced bool   `json:"forced"`
}

type PullReqCreatedPayload struct {
	BaseSegment
	PullReqSegment
	PullReqTargetReferenceSegment
	ReferenceSegment
	ReferenceDetailsSegment
}

type PullReqSegment struct {
	PullReq PullReqInfo `json:"pull_req"`
}

type PullReqInfo struct {
	Number        int64         `json:"number"`
	State         string        `json:"state"`
	IsDraft       bool          `json:"is_draft"`
	Title         string        `json:"title"`
	Description   string        `json:"description"`
	SourceRepoID  int64         `json:"source_repo_id"`
	SourceBranch  string        `json:"source_branch"`
	TargetRepoID  int64         `json:"target_repo_id"`
	TargetBranch  string        `json:"target_branch"`
	MergeStrategy *string       `json:"merge_strategy,omitempty"`
	Author        PrincipalInfo `json:"author"`
	PrURL         string        `json:"pr_url"`
}

type PullReqTargetReferenceSegment struct {
	TargetRef ReferenceInfo `json:"target_ref"`
}

type PullReqReopenedPayload PullReqCreatedPayload

type PullReqBranchUpdatedPayload struct {
	BaseSegment
	PullReqSegment
	PullReqTargetReferenceSegment
	ReferenceSegment
	ReferenceDetailsSegment
	ReferenceUpdateSegment
}

type PullReqMergedPayload PullReqClosedPayload

type PullReqClosedPayload struct {
	BaseSegment
	PullReqSegment
	PullReqTargetReferenceSegment
	ReferenceSegment
	ReferenceDetailsSegment
}

type PullReqCommentPayload struct {
	BaseSegment
	PullReqSegment
	PullReqTargetReferenceSegment
	ReferenceSegment
	ReferenceDetailsSegment
	PullReqCommentSegment
}

type PullReqCommentSegment struct {
	CommentInfo CommentInfo `json:"comment"`
}

type CommentInfo struct {
	ID       int64  `json:"id"`
	ParentID *int64 `json:"parent_id,omitempty"`
	Text     string `json:"text"`
}

type PullReqReviewerCreatedPayload PullReqReviewerChangedPayload
type PullReqReviewerDeletedPayload PullReqReviewerChangedPayload

type PullReqReviewerChangedPayload struct {
	BaseSegment
	PullReqSegment
	ReviewerSegment
}

type ReviewerSegment struct {
	Reviewer PrincipalInfo `json:"reviewer"`
}

type PullReqReviewSubmittedPayload struct {
	BaseSegment
	Author   *PrincipalInfo
	Reviewer *PrincipalInfo
	Decision string
}
