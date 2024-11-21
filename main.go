package main

import (
	"fmt"
	"freb/cmd"
	"freb/config"
	"freb/utils"
	"os"
	"runtime/pprof"
)

func main() {
	f, _ := os.OpenFile("cpu.profile", os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	err := config.GetConfig()
	config.Cfg.TmpDir = os.TempDir()
	if err != nil {
		utils.Err(fmt.Errorf("配置错误: %v", err))
	}
	cmd.Execute()
}
