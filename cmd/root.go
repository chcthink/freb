package cmd

import (
	"errors"
	"fmt"
	"freb/config"
	"freb/models"
	"freb/source"
	"freb/source/sources"
	"freb/utils"
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
	rootCmd.PersistentFlags().StringVarP(&novel.Format, formatCmd, "f", "epub", "转换至指定格式 默认 epub")
	rootCmd.PersistentFlags().StringVarP(&novel.Id, idCmd, "i", "", "下载书本id")
	rootCmd.PersistentFlags().StringVarP(&novel.Cover, coverCmd, "c", coverDefault, "封面路径")
	rootCmd.PersistentFlags().StringVarP(&novel.Out, outCmd, "o", "", "输出文件名")
	rootCmd.PersistentFlags().StringVarP(&novel.ContentImg, subCmd, "s", contentImgDefault, "每章标题logo")
	rootCmd.PersistentFlags().StringVarP(&novel.IntroImg, descImgCmd, "e", introImgDefault, "内容介绍logo")
	rootCmd.PersistentFlags().StringVarP(&novel.Vol, volCmd, "b", volDefault, "卷logo")
	rootCmd.PersistentFlags().StringVarP(&novel.Author, authorCmd, "a", "Unknown", "作者")
	rootCmd.PersistentFlags().BoolVarP(&novel.Desc, descCmd, "d", true, "是否包含制作说明,默认包含,使用 -d 来取消包含")
	rootCmd.PersistentFlags().StringVarP(&novel.Lang, langCmd, "l", "zh-Hans", "默认中文zh-Hans,英文 en")
	rootCmd.PersistentFlags().StringVarP(&novel.Path, pathCmd, "p", "", "转化txt路径")
	// rootCmd.PersistentFlags().StringVarP(&novel.Font, fontCmd, "w", "", "正文字体字体")
}

var novel models.Book

var rootCmd = &cobra.Command{
	Use:   "freb",
	Short: "freb用于下载小说并转换至指定格式",
	Run: func(cmd *cobra.Command, args []string) {
		err := BookParamCheck()
		if err != nil {
			utils.Err(err)
			return
		}
		if cmd.PersistentFlags().Changed(descCmd) {
			novel.Desc = false
		}
		var source source.Source
		if len(novel.Id) > 0 {
			source = &sources.UrlSource{}
		}
		path := novel.Path
		if len(args) > 0 {
			path = args[0]
		}
		if len(path) > 0 {
			source = &sources.TxtSource{}
			novel.Name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			novel.Name, novel.Author = utils.GetBookInfo(novel.Name)
		}
		err = source.GetBook(&novel)
		if err != nil {
			utils.Err(err)
			return
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.Err(err)
		os.Exit(1)
	}
}

func BookParamCheck() (err error) {
	if len(novel.Path) == 0 && len(novel.Id) == 0 {
		return errors.New(sourceErr)
	}
	if len(novel.Path) > 0 {
		if !utils.IsFileExist(novel.Path) {
			return errors.New(fmt.Sprintf(sourceErr, novel.Path))
		}
	}

	if len(novel.Id) > 0 {
		novel.IsOld = utils.CheckNum(novel.Id)
	}

	// cover subCover vol
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
		return err
	}
	novel.IntroImg, err = utils.SetImage(novel.IntroImg, config.Cfg.TmpDir, introImgDefault, nil)
	if err != nil {
		return err
	}
	novel.ContentImg, err = utils.SetImage(novel.ContentImg, config.Cfg.TmpDir, contentImgDefault, nil)
	if err != nil {
		return err
	}
	novel.Vol, err = utils.SetImage(novel.Vol, config.Cfg.TmpDir, volDefault, nil)
	if err != nil {
		return err
	}
	return nil
}
