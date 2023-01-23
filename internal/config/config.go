package config

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"

	"gihub.com/wanyuqin/gtp/internal/object"
)

type Proto struct {
	BasicSet BasicSet         `json:"basic_set"`
	Message  []object.Message `json:"message"`
}

type BasicSet struct {
	PackageName   string `json:"package_name" yaml:"basicSet.packageName"`
	GoPackageName string `json:"go_package_name" yaml:"goPackageName"`
	ApiVersion    string `json:"api_version" yaml:"apiVersion"`
	OutputPath    string `json:"output_path" yaml:"outputPath"`
	FileName      string `json:"file_name"`
}

var (
	defaultPath       = "./configs"
	defaultFileName   = "default-config.yaml"
	defaultConfigType = "yaml"
)

func InitConfig(configPath string) (BasicSet, error) {
	b := BasicSet{}
	if configPath == "" {
		configPath = path.Join(defaultPath, defaultFileName)
	}

	fi, err := os.Stat(configPath)
	if err != nil {
		return b, err
	}

	fmt.Println(fi.Name())

	viper.SetConfigType(defaultConfigType)
	// viper.AddConfigPath(configPath)
	viper.SetConfigFile(configPath)
	// viper.SetConfigName(fi.Name())
	err = viper.ReadInConfig()
	if err != nil {
		return b, err
	}
	err = viper.Unmarshal(&b)

	if err != nil {
		return b, err
	}

	return b, nil
}
