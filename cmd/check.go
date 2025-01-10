package cmd

import (
	"freb/models"
	"freb/source/sources"
	"freb/utils/stdout"
	"github.com/spf13/cobra"
)

const (
	testBaseInfoCmd = "chapter"
)

var (
	testBase string
)

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().StringVar(&testBase, testBaseInfoCmd, "", "测试章节内容,值为章节 url")
}

var testCmd = &cobra.Command{
	Use:   "check",
	Short: "测试 URL 获取书籍时,正则与 xpath 是否匹配",
	Run: func(cmd *cobra.Command, args []string) {
		catch, err := checkTestFlag()
		if err != nil {
			stdout.Errln(err)
			return
		}
		content, err := sources.SetSection(testBase, "", &ef, catch)
		if err != nil {
			stdout.Errln(err)
			return
		}
		stdout.Contentln(content)
	},
}

func checkTestFlag() (*models.BookCatch, error) {
	return CheckCatchUrl(testBase)
}
