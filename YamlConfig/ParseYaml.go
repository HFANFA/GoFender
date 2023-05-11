package YamlConfig

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	KafkaServer       string `yaml:"kafka_server"`
	MysqlServer       string `yaml:"mysql_dns"`
	MysqlUser         string `yaml:"mysql_user"`
	MysqlPassword     string `yaml:"mysql_password"`
	MysqlDatabase     string `yaml:"mysql_database"`
	MysqlMaxIdleConns int    `yaml:"mysql_max_idle_conns"`
	MysqlMaxOpenConns int    `yaml:"mysql_max_open_conns"`
	WebAddr           string `yaml:"web_addr"`
	RulesPath         string `yaml:"rules_path"`
	IpToLocationDb    string `yaml:"ip_location_db"`
	ModelPath         string `yaml:"model_path"`
	TempPath          string `yaml:"temp_path"`
	YaraFile          string `yaml:"yara_file"`
}

var Myconfig Config

func ParseYaml(YamlFilePath string) Config {
	var config Config
	configFile, err := os.ReadFile(YamlFilePath)
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal("read config file error,the reason as follow:", err)
	}
	return config
}
