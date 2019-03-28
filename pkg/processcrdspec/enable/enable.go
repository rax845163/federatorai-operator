package enable

import (
	"github.com/containers-ai/federatorai-operator/pkg/utils"
)

func IgnoreGUIYAML(fileLocation []string, guicomponent []string) []string {
	for _, v := range guicomponent {
		index := utils.IndexOf(v, fileLocation)
		if index != -1 {
			fileLocation = append(fileLocation[:index], fileLocation[index+1:]...)
		}
	}
	return fileLocation
}
func IgnoreExcutionYAML(fileLocation []string, Excutioncomponent []string) []string {
	for _, v := range Excutioncomponent {
		index := utils.IndexOf(v, fileLocation)
		if index != -1 {
			fileLocation = append(fileLocation[:index], fileLocation[index+1:]...)
		}
	}
	return fileLocation
}
