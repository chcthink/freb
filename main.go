package main

import (
	"fmt"
	"freb/cmd"
	"freb/config"
	"freb/utils"
)

func main() {
	err := config.GetConfig()
	if err != nil {
		utils.Err(fmt.Errorf("配置错误: %v", err))
	}
	cmd.Execute()
}
