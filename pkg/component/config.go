package component

type InfluxDBConfig struct {
	Address  string
	Username string
	Password string
}

type PrometheusConfig struct {
	Address  string
	Username string
	Password string
}

type FederatoraiAgentGPUDatasourceConfig struct {
	InfluxDB   InfluxDBConfig
	Prometheus PrometheusConfig
}

type FederatoraiAgentGPUConfig struct {
	Datasource FederatoraiAgentGPUDatasourceConfig
}

func NewDefaultFederatoraiAgentGPUConfig() FederatoraiAgentGPUConfig {
	return FederatoraiAgentGPUConfig{
		Datasource: FederatoraiAgentGPUDatasourceConfig{
			InfluxDB: InfluxDBConfig{
				Address:  "",
				Username: "",
				Password: "",
			},
			Prometheus: PrometheusConfig{
				Address:  "",
				Username: "",
				Password: "",
			},
		},
	}
}
