package snailjob

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type SysConfig struct {
	SnailJob struct {
		Server struct {
			Port int    `yaml:"port"`
			Host string `yaml:"host"`
		} `yaml:"server"`
		Namespace string `yaml:"namespace"`
		Group     string `yaml:"group"`
		Token     string `yaml:"token"`
		Host      string `yaml:"host"`
		Port      string `yaml:"port"`
	} `yaml:"snail-job"`
}

func NewConfig() *SysConfig {
	// 读取 YAML 文件
	file, err := os.Open("config.yml")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// 解码 YAML 文件
	var config SysConfig
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	return &config

}
