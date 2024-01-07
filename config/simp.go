package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type SimpConfig struct {
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

func NewConfig(path string) (conf SimpConfig, err error) {
	wd, _ := os.Getwd()
	configPath := filepath.Join(wd, path)
	// 读取 YAML 文件
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		return conf, err
	}

	// 解析 YAML 数据
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		return conf, err
	}

	// 打印解析后的配置
	fmt.Printf("%+v\n", conf)
	return conf, nil
}

func ResetConfig(yamlContent string, filePath string) error {
	// 写入 YAML 内容到文件
	err := os.WriteFile(filePath, []byte(yamlContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing YAML to file: %v", err)
	}
	return nil
}
