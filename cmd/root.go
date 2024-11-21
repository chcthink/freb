package cmd

import (
	"freb/models"
	"freb/source"
	"freb/utils"
	"github.com/spf13/cobra"
	"os"
)

var (
	isDownload bool // 书籍来源
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&Novel.Format, "format", "f", "epub", "转换至指定格式 默认 epub")
	rootCmd.PersistentFlags().StringVarP(&Novel.Url, "url", "u", "", "下载网址 不能为空")
	rootCmd.PersistentFlags().StringVarP(&Novel.Cover, "cover", "c", "", "封面路径")
	rootCmd.PersistentFlags().StringVarP(&Novel.Out, "out", "o", "", "输出文件名")
	rootCmd.PersistentFlags().StringVarP(&Novel.SubCover, "sub", "s", "", "每章封面")
	rootCmd.PersistentFlags().StringVarP(&Novel.Author, "author", "a", "Unknown", "作者")
	rootCmd.PersistentFlags().StringVarP(&Novel.Author, "font", "w", "assets/fonts/975MaruSC-Medium.ttf", "作者")
	rootCmd.PersistentFlags().BoolVarP(&isDownload, " is_download", "d", true, "书籍来源:下载-true")
}

var Novel models.Book

var rootCmd = &cobra.Command{
	Use:   "freb",
	Short: "freb用于下载小说并转换至指定格式",
	Run: func(cmd *cobra.Command, args []string) {
		switch isDownload {
		case true:
			var source source.UrlSource
			err := source.GetBook(&Novel)
			if err != nil {
				utils.Err(err)
				return
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.Err(err)
		os.Exit(1)
	}
}
