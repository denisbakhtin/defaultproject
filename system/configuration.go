package system

import (
	"encoding/json"
	"fmt"
)

type ConfigurationDatabase struct {
	Host     string `json:"host"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Configurations struct {
	Development Configuration
	Production  Configuration
	Testing     Configuration
}
type Configuration struct {
	Secret   string `json:"secret"`
	Database ConfigurationDatabase
}

func LoadConfiguration(env *string, data []byte) (*Configuration, error) {
	confs := &Configurations{}
	err := json.Unmarshal(data, &confs)
	switch *env {
	case "production":
		return &confs.Production, err
	case "development":
		return &confs.Development, err
	case "testing":
		return &confs.Testing, err
	default:
		return nil, fmt.Errorf("Unknown environment flag")
	}
}
