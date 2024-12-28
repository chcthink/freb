package config

import (
	"fmt"
	"freb/utils"
	"github.com/pelletier/go-toml"
	"os"
)

type Page struct {
	Title string
	Dom   string
}

type Style struct {
	Instruction Page
	Desc        Page
	Vol         string
	Chapter     string
}

type Selector struct {
	Title  string
	Author string
	Intro  string
}

type Config struct {
	*Style
	*Selector
	TmpDir    string `toml:"-"`
	DelayTime int    `toml:"delay_time"`
	Cookies   map[string]string
}

var Cfg Config

const (
	cfgPath = "config.toml"
	cfgErr  = "初始化配置文件错误: %s"
)

func GetConfig() error {
	Cfg.TmpDir = os.TempDir()
	source, err := utils.LocalOrDownload(cfgPath, Cfg.TmpDir)
	if err != nil {
		return fmt.Errorf(cfgErr, err)
	}
	file, err := toml.LoadFile(source)
	if err != nil {
		return fmt.Errorf(cfgErr, err)
	}
	err = file.Unmarshal(&Cfg)
	if err != nil {
		return fmt.Errorf(cfgErr, err)
	}
	return nil
}
