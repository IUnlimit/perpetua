package model

import "time"

type Workflow struct {
	TotalCount   int            `json:"total_count"`
	WorkflowRuns []*WorkflowRun `json:"workflow_runs"`
}

type WorkflowRun struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	HeadSha      string    `json:"head_sha"`
	Path         string    `json:"path"`
	DisplayTitle string    `json:"display_title"`
	Status       string    `json:"status"`
	Conclusion   string    `json:"conclusion"`
	WorkflowID   int       `json:"workflow_id"`
	URL          string    `json:"url"`
	HTMLURL      string    `json:"html_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	JobsURL      string    `json:"jobs_url"`
	ArtifactsURL string    `json:"artifacts_url"`
}

type Artifact struct {
	TotalCount int         `json:"total_count"`
	Artifacts  []Artifacts `json:"artifacts"`
}

type Artifacts struct {
	ID                 int       `json:"id"`
	NodeID             string    `json:"node_id"`
	Name               string    `json:"name"`
	SizeInBytes        int       `json:"size_in_bytes"`
	URL                string    `json:"url"`
	ArchiveDownloadURL string    `json:"archive_download_url"`
	Expired            bool      `json:"expired"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	ExpiresAt          time.Time `json:"expires_at"`
}
