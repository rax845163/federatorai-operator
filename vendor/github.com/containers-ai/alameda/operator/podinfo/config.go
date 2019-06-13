package podinfo

type Config struct {
	LabelsFile string
}

func NewConfig() *Config {
	return &Config{
		LabelsFile: "/etc/podinfo/labels",
	}
}
