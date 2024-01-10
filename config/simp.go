package config

import (
	"Simp/utils"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
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
	SIMP_PRODUCTION := os.Getenv("SIMP_PRODUCTION")
	fmt.Println("SIMP_PRODUCTION", SIMP_PRODUCTION)
	var configPath string
	if SIMP_PRODUCTION == "Yes" {
		SIMP_CONFIG_PATH := os.Getenv("SIMP_CONFIG_PATH")
		fmt.Println("SIMP_CONFIG_PATH", SIMP_CONFIG_PATH)
		configPath = SIMP_CONFIG_PATH
	} else {
		configPath = filepath.Join(wd, path)
	}
	// 读取 YAML 文件
	yamlFile, err := os.ReadFile(configPath)
	fmt.Println("Get FilePath from ", configPath)
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

func CoverConfig(config SimpConfig, filePath string) error {
	// 删除
	utils.IFExistThenRemove(filePath)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	if err := encoder.Encode(config); err != nil {
		return err
	}

	return nil
}

func ParseConfig(yamlString string) (SimpConfig, error) {
	var config SimpConfig
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		return SimpConfig{}, err
	}
	return config, nil
}
