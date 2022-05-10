package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var instance *Config = nil

func Get() *Config {
	if instance != nil {
		return instance
	}

	config := &Config{}
	path := "config.json"

	j, err := os.ReadFile(path)
	if err != nil {
		panic("config not found")
	}

	err = json.Unmarshal(j, config)
	if err != nil {
		panic("config error")
	}

	// 处理字体文件
	config.FontMap = map[string]string{}
	for _, font := range config.Fonts {
		if strings.Index(font.Path, config.Roots.Fonts) == 0 {
			config.FontMap[font.Family] = font.Path
		} else {
			config.FontMap[font.Family] = fmt.Sprintf("%s/%s", config.Roots.Fonts, font.Path)
		}
	}

	instance = config

	return config
}
