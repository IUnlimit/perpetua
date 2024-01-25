package hook

import (
	"errors"
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/IUnlimit/perpetua/internal/utils"
)

const TOKEN = "0272493805f67556a6150d6576c161970a5f3a8bf0e8a56cd4aca4bac9a2b7b7b078878a34ed81147981e5268ab0462d64749bdb5d2fb3b0ca3c13bab71ed410"

const ActionRunsApi = "https://api.github.com/repos/{0}/{1}/actions/runs"
const ActionArtifactsApi = "https://api.github.com/repos/{0}/{1}/actions/artifacts"

var headers = map[string]string{
	"Accept":               "application/vnd.github+json",
	"X-GitHub-Api-Version": "2022-11-28",
}

// GetWorkflowRunIds 获取 owner/repo 仓库的指定 name 下的 workflow action id
func GetWorkflowRunIds(owner string, repo string, params map[string]string, name string) ([]int64, error) {
	url := utils.Format(ActionRunsApi, owner, repo)
	url, err := utils.BuildURLParams(url, params)
	if err != nil {
		return nil, err
	}

	var workflow model.WorkflowAction
	err = utils.GetJson(url, headers, &workflow)
	if err != nil {
		return nil, err
	}
	if workflow.TotalCount == 0 {
		return nil, errors.New("can't match any workflow")
	}

	var ids []int64
	for _, run := range workflow.WorkflowRuns {
		if run.Name == name {
			ids = append(ids, run.ID)
		}
	}
	return ids, nil
}

// GetArtifactsUrls 获取 owner/repo 仓库的产出工件
func GetArtifactsUrls(owner string, repo string, params map[string]string) ([]*model.Artifact, error) {
	url := utils.Format(ActionArtifactsApi, owner, repo)
	url, err := utils.BuildURLParams(url, params)
	if err != nil {
		return nil, err
	}

	var artifact model.ArtifactAction
	err = utils.GetJson(url, headers, &artifact)
	if err != nil {
		return nil, err
	}
	if artifact.TotalCount == 0 {
		return nil, errors.New("can't match any artifact")
	}

	return artifact.Artifacts, nil
}

func GetAuthorizedFile(url string, filePath string, fileSize int64) error {
	err := utils.DownloadFileWithHeaders(url, filePath, headers, fileSize)
	if err != nil {
		return err
	}
	return nil
}
