package main

import (
	"fmt"
	"freb/cmd"
	"freb/config"
	"freb/utils/reg"
	"freb/utils/stdout"
	"os"
)

func main() {
	err := config.InitConfig()
	defer os.RemoveAll(config.Cfg.TmpDir)
	reg.InitCustomMatchReg(config.Cfg.Regs)
	if err != nil {
		stdout.Errln(fmt.Errorf("配置错误: %v", err))
		return
	}
	cmd.Execute()
}
