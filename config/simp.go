package config

import (
	"Simp/utils"
	"fmt"
	"os"
	"path/filepath"

	jsonToYaml "github.com/ghodss/yaml"
	"gopkg.in/yaml.v2"
)

type SimpConfig struct {
	Server struct {
		Name       string `yaml:"name" `
		Port       int    `yaml:"port" `
		Type       string `yaml:"type"`
		StaticPath string `yaml:"staticPath" `
		Storage    string `yaml:"storage" `
		Main       bool   `yaml:"main" `
		Proxy      []struct {
			Server struct {
				Type string `yaml:"type"`
				Name string `yaml:"name" `
			} `yaml:"server"`
		} `yaml:"proxy"`
	} `yaml:"server"`
}

type CoverConfigVo struct {
	Conf       SimpConfig
	ServerName string
}

func NewConfig(path string) (conf SimpConfig, err error) {
	wd, _ := os.Getwd()
	SIMP_PRODUCTION := os.Getenv("SIMP_PRODUCTION")
	fmt.Println("NewConfig SIMP_PRODUCTION", SIMP_PRODUCTION)
	var configPath string
	if SIMP_PRODUCTION == "Yes" {
		SIMP_SERVER_PATH := os.Getenv("SIMP_SERVER_PATH")
		prodConfigPath := filepath.Join(SIMP_SERVER_PATH, "simpProd.yaml")
		configPath = filepath.Join(SIMP_SERVER_PATH, "simp.yaml")
		yamlFile, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Println("Error to Get YamlFile ", configPath)
		}
		yamlProdFile, err := os.ReadFile(prodConfigPath)
		if err != nil {
			fmt.Println("Error to Get yamlProdFile ", prodConfigPath)
		}
		var doc1, doc2 SimpConfig
		if err := yaml.Unmarshal(yamlProdFile, &doc1); err != nil {
			fmt.Println("无法解析第一个YAML文件：", err)
		}
		if err := yaml.Unmarshal(yamlFile, &doc2); err != nil {
			fmt.Println("无法解析第二个YAML文件：", err)
		}
		fmt.Println("doc1", doc1)
		fmt.Println("doc2", doc2)
		mergeYaml := MergeYAML(doc1, doc2)
		// 将 interface{} 转换为 MyStruct
		if conf, ok := mergeYaml.(SimpConfig); ok {
			return conf, nil
			// 转换成功，可以使用 myStruct 进行操作
		} else {
			fmt.Println("Error to Parsse Interface to Yaml ", prodConfigPath)
			return conf, err
			// 转换失败，处理错误情况
		}
	} else {
		configPath = filepath.Join(wd, path)
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
}

func ResetConfig(yamlContent string, filePath string) error {
	// 写入 YAML 内容到文件
	err := os.WriteFile(filePath, []byte(yamlContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing YAML to file: %v", err)
	}
	return nil
}

func CoverConfig(content string, filePath string) error {
	// 删除
	utils.IFExistThenRemove(filePath)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

func ParseConfig(yamlString string) (string, error) {
	yml, err := jsonToYaml.JSONToYAML([]byte(yamlString))
	if err != nil {
		fmt.Println("JSON TO YamlError")
	}
	fmt.Println("Cover yml \n", string(yml))

	return string(yml), nil
}

// mergeYAML 递归合并两个YAML文档
func MergeYAML(doc1, doc2 interface{}) interface{} {
	switch doc1 := doc1.(type) {
	case map[interface{}]interface{}:
		doc2, ok := doc2.(map[interface{}]interface{})
		if !ok {
			return doc1
		}
		merged := make(map[interface{}]interface{})
		for k, v := range doc1 {
			merged[k] = MergeYAML(v, doc2[k])
		}
		return merged
	default:
		return doc1
	}
}
