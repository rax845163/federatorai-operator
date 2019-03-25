package enable

import (
	"github.com/containers-ai/federatorai-operator/pkg/utils"
)

func IgnoreGUIYAML(FileLocation []string, Guicomponent []string) []string {
	for _, v := range Guicomponent {
		index := utils.IndexOf(v, FileLocation)
		if index != -1 {
			FileLocation = append(FileLocation[:index], FileLocation[index+1:]...)
		}
	}
	return FileLocation
}
func IgnoreExcutionYAML(FileLocation []string, Excutioncomponent []string) []string {
	for _, v := range Excutioncomponent {
		index := utils.IndexOf(v, FileLocation)
		if index != -1 {
			FileLocation = append(FileLocation[:index], FileLocation[index+1:]...)
		}
	}
	return FileLocation
}
