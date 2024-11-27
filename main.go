package main

import (
	"fmt"
	"freb/cmd"
	"freb/config"
	"freb/utils"
	"os"
)

func main() {
	err := config.GetConfig()
	config.Cfg.TmpDir = os.TempDir()
	if err != nil {
		utils.Err(fmt.Errorf("配置错误: %v", err))
	}
	cmd.Execute()
}
