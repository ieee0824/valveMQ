package valve

import (
	"encoding/json"
	"os"
)

type dbConfig struct {
	DbUser string `json:"db_user"`
	DbPass string `json:"db_pass"`
	DbHost string `json:"db_host"`
	DbPort string `json:"db_port"`
	DbName string `json:"db_name"`
}

func (d *dbConfig) new() {
	if d.DbUser == "" {
		d.DbUser = "root"
	}
	if d.DbHost == "" {
		d.DbHost = "127.0.0.1"
	}
	if d.DbPort == "" {
		d.DbPort = "3306"
	}
	if d.DbName == "" {
		d.DbName = "mq"
	}
}

type apiConfig struct {
	APIPort string `json:"api_port"`
}

func (a *apiConfig) new() {
	if a.APIPort == "" {
		a.APIPort = "8080"
	}
}

type queueConfig struct {
	DequeueLimit uint `json:"dequeue_limit"`
}

func (q *queueConfig) new() {
}

type Config struct {
	dbConfig
	apiConfig
	queueConfig
}

func NewConfig(cfgFileName ...string) *Config {
	ret := &Config{}
	for _, name := range cfgFileName {
		f, err := os.Open(name)
		if err == nil {
			json.NewDecoder(f).Decode(ret)
			f.Close()
			break
		}
	}
	ret.dbConfig.new()
	ret.apiConfig.new()
	ret.queueConfig.new()

	return ret
}

func GetConfigPath() []string {
	return []string{
		"/etc/vmq.cnf",
		"./vmq.cnf",
	}
}
