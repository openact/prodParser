package conf

import (
	"github.com/openact/lib/utils"
)

var Runs []*Run
var DataPaths = utils.Conf.GetStringSlice("inputPaths")
var OutputPath = utils.Conf.GetString("outputPath")

func init() {
	Runs = getRunsInfo()
}

// define run params
type Run struct {
	Name     string   `yaml:"name"`
	Folder   string   `yaml:"folder"`
	CodeRepo []string `yaml:"codeRepo"`
}

func getRunsInfo() []*Run {
	var runs []*Run
	err := utils.Conf.UnmarshalKey("runs", &runs)
	if err != nil {
		panic(err)
	}

	for _, run := range runs {
		for i, repo := range run.CodeRepo {
			run.CodeRepo[i] = utils.GetFilePath(repo, DataPaths)
		}
	}
	return runs
}
