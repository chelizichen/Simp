package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type OriginConfig struct {
	Server struct {
		Name       string `yaml:"name"`
		Port       int    `yaml:"port"`
		Type       string `yaml:"type"`
		StaticPath string `yaml:"staticPath"`
		Storage    string `yaml:"storage"`
		Proxy      []struct {
			Server struct {
				Type string `yaml:"type"`
				Name string `yaml:"name"`
			} `yaml:"server"`
		} `yaml:"proxy"`
	} `yaml:"server"`
}

func main() {
	wd, _ := os.Getwd()
	config := wd + "/test/config/simp.yaml"
	// 读取 YAML 文件
	yamlFile, err := os.ReadFile(config)
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		return
	}

	// 解析 YAML 数据
	var conf OriginConfig
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		return
	}

	// 打印解析后的配置
	fmt.Printf("%+v\n", conf)
}
