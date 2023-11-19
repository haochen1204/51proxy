package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Url51      string `yaml:"51url"`
	FofaEmail  string `yaml:"fofaEmail"`
	FofaApiKey string `yaml:"fofaApiKey"`
}

func ReadConfig() *Config {
	// 读取 YAML 文件内容
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("无法读取文件: %v", err)
	}

	// 解析 YAML 文件内容
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("无法解析 YAML 文件: %v", err)
	}
	return &config
}
