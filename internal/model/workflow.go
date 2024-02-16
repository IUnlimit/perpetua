package model

import "time"

type WorkflowAction struct {
	TotalCount   int64          `json:"total_count"`
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
	WorkflowID   int64     `json:"workflow_id"`
	URL          string    `json:"url"`
	HTMLURL      string    `json:"html_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	JobsURL      string    `json:"jobs_url"`
	ArtifactsURL string    `json:"artifacts_url"`
}

type ArtifactAction struct {
	TotalCount int64       `json:"total_count"`
	Artifacts  []*Artifact `json:"artifacts"`
}

type Artifact struct {
	ID                 int64        `json:"id"`
	NodeID             string       `json:"node_id"`
	Name               string       `json:"name"`
	SizeInBytes        int64        `json:"size_in_bytes"`
	URL                string       `json:"url"`
	ArchiveDownloadURL string       `json:"archive_download_url"`
	Expired            bool         `json:"expired"`
	CreatedAt          time.Time    `json:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at"`
	ExpiresAt          time.Time    `json:"expires_at"`
	WorkflowRun        *WorkflowRun `json:"workflow_run"`
}
