package config

import (
	"fmt"
	"freb/models"
	"freb/utils"
	"freb/utils/stdout"
	"github.com/pelletier/go-toml"
	"net/http"
	"os"
)

var Cfg models.Config

const (
	cfgPath = "config.toml"
	cfgErr  = "初始化配置文件错误: %s"
)

func InitConfig() (err error) {
	initConfig()
	Cfg.TmpDir = os.TempDir()
	var source string
	if Cfg.From != "" {
		stdout.Fmtfln("正在从远程仓库下载文件: %s", Cfg.From)
		source, err = utils.DownloadTmp(Cfg.TmpDir, cfgPath, func() *http.Request {
			return utils.GetWithUserAgent(Cfg.From)
		})
	} else {
		source, err = utils.LocalOrDownload(cfgPath, Cfg.TmpDir, Cfg.From)
	}
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

func initConfig() {
	Cfg.BookCatch = make(map[string]*models.BookCatch)
	Cfg.InfoSelector = make(map[string]*models.InfoSelector)
}
