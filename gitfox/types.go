package gitfox

import "time"

type Repo struct {
	ID            int    `json:"id"`
	Path          string `json:"path"`
	Identifier    string `json:"identifier"`
	DefaultBranch string `json:"default_branch"`
	GitURL        string `json:"git_url"`
	UID           string `json:"uid"`
}

type User struct {
	ID          int    `json:"id"`
	UID         string `json:"uid"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Type        string `json:"type"`
	Created     int64  `json:"created"`
	Updated     int64  `json:"updated"`
}

type Identity struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CommitUser struct {
	Identity Identity  `json:"identity"`
	When     time.Time `json:"when"`
}

type Commit struct {
	Sha       string     `json:"sha"`
	Message   string     `json:"message"`
	Author    CommitUser `json:"author"`
	Committer CommitUser `json:"committer"`
	Added     []any      `json:"added"`
	Removed   []any      `json:"removed"`
	Modified  []string   `json:"modified"`
}

type Ref struct {
	Name string `json:"name"`
	Repo Repo   `json:"repo"`
}

type BranchUpdatedPayload struct {
	Trigger           string   `json:"trigger"`
	Repo              Repo     `json:"repo"`
	Principal         User     `json:"principal"`
	Ref               Ref      `json:"ref"`
	Sha               string   `json:"sha"`
	HeadCommit        Commit   `json:"head_commit"`
	Commits           []Commit `json:"commits"`
	TotalCommitsCount int      `json:"total_commits_count"`
	Commit            Commit   `json:"commit"`
	OldSha            string   `json:"old_sha"`
	Forced            bool     `json:"forced"`
}

type BranchCreatedPayload struct {
	Trigger    string `json:"trigger"`
	Repo       Repo   `json:"repo"`
	Principal  User   `json:"principal"`
	Ref        Ref    `json:"ref"`
	Sha        string `json:"sha"`
	HeadCommit Commit `json:"head_commit"`
	Commit     Commit `json:"commit"`
	OldSha     string `json:"old_sha"`
	Forced     bool   `json:"forced"`
}

type BranchDeletedPayload struct {
	Trigger   string `json:"trigger"`
	Repo      Repo   `json:"repo"`
	Principal User   `json:"principal"`
	Ref       Ref    `json:"ref"`
	Sha       string `json:"sha"`
	OldSha    string `json:"old_sha"`
	Forced    bool   `json:"forced"`
}

type TagCreatedPayload struct {
	Trigger    string `json:"trigger"`
	Repo       Repo   `json:"repo"`
	Principal  User   `json:"principal"`
	Ref        Ref    `json:"ref"`
	Sha        string `json:"sha"`
	HeadCommit Commit `json:"head_commit"`
	Commit     Commit `json:"commit"`
	OldSha     string `json:"old_sha"`
	Forced     bool   `json:"forced"`
}
