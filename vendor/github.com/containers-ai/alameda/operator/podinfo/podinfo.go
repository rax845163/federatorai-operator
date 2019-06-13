package podinfo

import (
	"io/ioutil"
	"regexp"
	"strings"
)

type PodInfo struct {
	Labels *map[string]string
}

func NewPodInfo(config *Config) *PodInfo {
	podInfo := &PodInfo{}
	podInfo.Labels = podInfo.parsePodLabels(config)
	return podInfo
}

func (*PodInfo) parsePodLabels(config *Config) *map[string]string {
	labelsFile := config.LabelsFile
	dat, _ := ioutil.ReadFile(labelsFile)
	labelsStr := strings.Replace(string(dat), "\n", " ", -1)
	rex := regexp.MustCompile("(.+?)=\"(.+?)\"")
	data := rex.FindAllStringSubmatch(labelsStr, -1)

	res := make(map[string]string)
	for _, kv := range data {
		k := strings.TrimSpace(kv[1])
		v := strings.TrimSpace(kv[2])
		res[k] = v
	}
	return &res
}
