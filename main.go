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
	if err != nil {
		utils.Err(err)
		return
	}
	config.Cfg.TmpDir = os.TempDir()
	if err != nil {
		utils.Err(fmt.Errorf("配置错误: %v", err))
	}
	cmd.Execute()
}
