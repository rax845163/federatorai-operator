package enable

import (
	"github.com/containers-ai/federatorai-operator/pkg/utils"
)

func DeleteGUIYAML(file_location []string, Guicomponent []string) []string {
	for _, v := range Guicomponent {
		index := utils.IndexOf(v, file_location)
		if index != -1 {
			file_location = append(file_location[:index], file_location[index+1:]...)
		}
	}
	return file_location
}
func DeleteExcutionYAML(file_location []string, Excutioncomponent []string) []string {
	for _, v := range Excutioncomponent {
		index := utils.IndexOf(v, file_location)
		if index != -1 {
			file_location = append(file_location[:index], file_location[index+1:]...)
		}
	}
	return file_location
}
