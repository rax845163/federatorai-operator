package assets

var (
	requiredConfigMaps = []string{
		"Requirement/ConfigMap/cluster-info.yaml",
	}
)

// GetRequiredConfigMaps returns configMap files that are required
func GetRequiredConfigMaps() []string {
	return requiredConfigMaps
}
