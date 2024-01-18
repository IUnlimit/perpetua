package perpetua

import (
	"fmt"
	"github.com/IUnlimit/perpetua/internal/logger"
	"github.com/IUnlimit/perpetua/tools"
)

func Login() {
	logger.Log.Println("Searching Lagrange.OneBot ...")
	initLagrange(false, "", 0)
}

func initLagrange(update bool, platform string, recordArtifactId int64) error {
	params := map[string]string{
		"per_page": "2",
		"status":   "success",
		"branch":   "master",
	}
	owner := "LagrangeDev"
	repo := "Lagrange.Core"

	ids, err := tools.GetWorkflowRunIds(owner, repo, params, "Lagrange.OneBot Build")
	if err != nil {
		return err
	}

	urlMap, err := tools.GetArtifactsUrls(owner, repo, ids[0])
	if err != nil {
		return err
	}

	logger.Log.Println("Please choose the Lagrange.OneBot software suitable for your platform (send the number before option)")
	var urls []string
	i := 0
	for name, url := range urlMap {
		urls = append(urls, url)
		fmt.Printf("[%d] %s\n", i, name)
		i++
	}

	var selectIndex int8
	fmt.Scanf("%d", &selectIndex)
	logger.Log.Println("Start downloading ...")
	fmt.Printf(urls[selectIndex])
	// TODO
	return nil
}
