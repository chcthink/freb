package main

import (
	"fmt"
	"freb/cmd"
	"freb/config"
	"freb/utils/stdout"
	"os"
	"path/filepath"
)

func main() {
	execPath, _ := os.Executable()
	_ = os.Chdir(filepath.Dir(execPath))
	err := config.GetConfig()
	if err != nil {
		stdout.Err(fmt.Errorf("配置错误: %v", err))
	}
	cmd.Execute()
}
