package enable

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k //index
		}
	}
	return -1 //not found.
}
func DeleteGUIYAML(file_location []string, Guicomponent []string) []string {
	for _, v := range Guicomponent {
		index := indexOf(v, file_location)
		if index != -1 {
			file_location = append(file_location[:index], file_location[index+1:]...)
		}
	}
	return file_location
}
func DeleteExcutionYAML(file_location []string, Excutioncomponent []string) []string {
	for _, v := range Excutioncomponent {
		index := indexOf(v, file_location)
		if index != -1 {
			file_location = append(file_location[:index], file_location[index+1:]...)
		}
	}
	return file_location
}
