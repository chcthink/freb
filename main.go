package main

import (
	"fmt"
	"freb/cmd"
	"freb/config"
	"freb/utils/stdout"
)

func main() {
	err := config.GetConfig()
	if err != nil {
		stdout.Err(fmt.Errorf("配置错误: %v", err))
	}
	cmd.Execute()
}
