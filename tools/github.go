package tools

import (
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/IUnlimit/perpetua/internal/utils"
	"strconv"
)

// TOKEN TODO 更换低权限token
const TOKEN = "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"

const ActionRunsApi = "https://api.github.com/repos/{0}/{1}/actions/runs"
const ActionArtifactsApi = "https://api.github.com/repos/{0}/{1}/actions/runs/{2}/artifacts"
const ArtifactDownloadApi = "https://github.com/{0}/{1}/actions/runs/{2}/artifacts/"

var headers = map[string]string{
	"Accept":               "application/vnd.github+json",
	"Authorization":        "Bearer " + TOKEN,
	"X-GitHub-Api-Version": "2022-11-28",
}

// GetWorkflowRunIds 获取 owner/repo 仓库的指定 name 下的 workflow action id
func GetWorkflowRunIds(owner string, repo string, params map[string]string, name string) ([]int64, error) {
	url := utils.Format(ActionRunsApi, owner, repo)
	url, err := utils.BuildURLParams(url, params)
	if err != nil {
		return nil, err
	}

	var workflow model.Workflow
	err = utils.GetJson(url, headers, &workflow)
	if err != nil {
		return nil, err
	}

	var ids []int64
	for _, run := range workflow.WorkflowRuns {
		if run.Name == name {
			ids = append(ids, run.ID)
		}
	}
	return ids, nil
}

func GetArtifactsUrls(owner string, repo string, workflowRunId int64) (map[string]string, error) {
	url := utils.Format(ActionArtifactsApi, owner, repo, workflowRunId)
	downloadUrl := utils.Format(ArtifactDownloadApi, owner, repo, workflowRunId)

	var artifact model.Artifact
	err := utils.GetJson(url, headers, &artifact)
	if err != nil {
		return nil, err
	}

	urls := make(map[string]string)
	for _, entity := range artifact.Artifacts {
		urls[entity.Name] = downloadUrl + strconv.Itoa(entity.ID)
	}
	return urls, nil
}
