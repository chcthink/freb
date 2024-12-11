package cmd

import (
	"errors"
	"fmt"
	"freb/config"
	"freb/formatter"
	"freb/models"
	"freb/source"
	"freb/source/sources"
	"freb/utils"
	"freb/utils/stdout"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	formatCmd  = "format"
	idCmd      = "id"
	coverCmd   = "cover"
	volCmd     = "vol"
	outCmd     = "out"
	subCmd     = "sub"
	descImgCmd = "img"
	authorCmd  = "author"
	descCmd    = "desc"
	langCmd    = "lang"
	pathCmd    = "path"
	// fontCmd   = "font"
)

const (
	coverDefault      = "cover.jpg"
	contentImgDefault = "content_logo.jpg"
	introImgDefault   = "intro_logo.jpg"
	volDefault        = "vol.jpg"
)

const (
	sourceErr = "文件路径或地址错误: %s"
)

func init() {
	// rootCmd.PersistentFlags().StringVarP(&novel.Format, formatCmd, "f", "epub", "转换至指定格式 默认 epub")
	rootCmd.PersistentFlags().StringVarP(&novel.Id, idCmd, "i", "", "下载书本id")
	rootCmd.PersistentFlags().StringVarP(&novel.Cover, coverCmd, "c", coverDefault, "封面路径")
	rootCmd.PersistentFlags().StringVarP(&novel.Out, outCmd, "o", "", "输出文件名")
	rootCmd.PersistentFlags().StringVarP(&novel.ContentImg, subCmd, "s", contentImgDefault, "每章标题logo")
	rootCmd.PersistentFlags().StringVarP(&novel.IntroImg, descImgCmd, "e", introImgDefault, "内容介绍logo")
	rootCmd.PersistentFlags().StringVarP(&novel.Vol, volCmd, "b", volDefault, "卷logo")
	rootCmd.PersistentFlags().StringVarP(&novel.Author, authorCmd, "a", "Unknown", "作者")
	rootCmd.PersistentFlags().BoolVarP(&novel.IsDesc, descCmd, "d", true, "是否包含制作说明,默认包含,使用 -d 来取消包含")
	rootCmd.PersistentFlags().StringVarP(&novel.Lang, langCmd, "l", "zh-Hans", "默认中文zh-Hans,英文 en")
	rootCmd.PersistentFlags().StringVarP(&novel.Path, pathCmd, "p", "", "转化txt路径")
	// rootCmd.PersistentFlags().StringVarP(&novel.Font, fontCmd, "w", "", "正文字体字体")
}

var novel models.Book

var rootCmd = &cobra.Command{
	Use:   "freb",
	Short: "freb用于下载小说并转换至EPub",
	Run: func(cmd *cobra.Command, args []string) {
		path := novel.Path
		if len(args) > 0 {
			path = args[0]
		}
		err := CheckFlag(cmd, path)
		if err != nil {
			stdout.Err(err)
			return
		}
		var source source.Source
		if len(novel.Id) > 0 {
			source = &sources.UrlSource{}
		}
		if len(path) > 0 {
			source = &sources.TxtSource{}
			novel.Name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			novel.Name, novel.Author = utils.GetBookInfo(novel.Name)
		}

		var ef formatter.EpubFormat
		ef.Book = &novel
		ef.AssetsPath = &formatter.AssetsPath{}
		err = InitAssets(ef)
		if err != nil {
			stdout.Err(err)
			return
		}

		err = source.GetBook(ef)
		if err != nil {
			stdout.Err(err)
			return
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		stdout.Err(err)
		os.Exit(1)
	}
}

func CheckFlag(cmd *cobra.Command, cmdPath string) (err error) {
	if cmd.PersistentFlags().Changed(descCmd) {
		novel.IsDesc = false
	}
	if len(cmdPath) == 0 && len(novel.Id) == 0 {
		return errors.New(sourceErr)
	}
	if len(cmdPath) > 0 {
		if !utils.IsFileInWorkDir(cmdPath) {
			return errors.New(fmt.Sprintf(sourceErr, cmdPath))
		}
	}

	if len(novel.Id) > 0 {
		novel.IsOld = utils.CheckNum(novel.Id)
	}
	return
}

func InitAssets(ef formatter.EpubFormat) (err error) {
	// image
	novel.Cover, err = utils.SetImage(novel.Cover, config.Cfg.TmpDir, coverDefault, func() *http.Request {
		var req *http.Request
		if len(novel.Id) > 0 {
			req = utils.NewGet(utils.CoverUrl(novel.IsOld, novel.Id))
			req.Header.Set("Referer", utils.SearchUrl(novel.IsOld))
			return req
		}
		return nil
	})
	if err != nil {
		return
	}
	novel.IntroImg, err = utils.SetImage(novel.IntroImg, config.Cfg.TmpDir, introImgDefault, nil)
	if err != nil {
		return
	}
	novel.ContentImg, err = utils.SetImage(novel.ContentImg, config.Cfg.TmpDir, contentImgDefault, nil)
	if err != nil {
		return
	}
	novel.Vol, err = utils.SetImage(novel.Vol, config.Cfg.TmpDir, volDefault, nil)
	if err != nil {
		return
	}
	// font
	ef.AssetsPath.Font, err = utils.LocalOrDownload("assets/fonts/font.ttf", config.Cfg.TmpDir)
	if err != nil {
		return
	}
	// MetaINF
	ef.AssetsPath.MetaInf, err = utils.LocalOrDownload("assets/META-INF/com.apple.ibooks.display-options.xml", config.Cfg.TmpDir)
	if err != nil {
		return
	}
	// css
	ef.AssetsPath.CommonCss, err = utils.LocalOrDownload("assets/styles/main.css", config.Cfg.TmpDir)
	if err != nil {
		return
	}
	ef.AssetsPath.FontCss, err = utils.LocalOrDownload("assets/styles/fonts.css", config.Cfg.TmpDir)
	if err != nil {
		return
	}
	ef.AssetsPath.CoverCss, err = utils.LocalOrDownload("assets/styles/cover.css", config.Cfg.TmpDir)
	if err != nil {
		return
	}
	ef.AssetsPath.InstructionCss, err = utils.LocalOrDownload("assets/styles/instruction.css", config.Cfg.TmpDir)
	if err != nil {
		return
	}
	return
}
