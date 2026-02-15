package config

type (
	Config struct {
		HTTP        `json:"http"`
		HealthCheck `json:"healthcheck"`
	}

	HTTP struct {
		Port string `json:"port"`
	}

	HealthCheck struct {
		Threads int `json:"threads"`
	}
)
