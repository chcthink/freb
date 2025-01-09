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
	"freb/utils/reg"
	"freb/utils/stdout"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	urlCmd        = "url"
	coverCmd      = "cover"
	volCmd        = "vol"
	outCmd        = "out"
	subCmd        = "sub"
	descImgCmd    = "img"
	authorCmd     = "author"
	descCmd       = "desc"
	langCmd       = "lang"
	pathCmd       = "path"
	jumpCmd       = "jump"
	delayCmd      = "delay"
	catalogUrlCmd = "curl"
	configCmd     = "config"
)

const (
	fromErr   = "数据来源异常(-p -i 皆为空)"
	sourceErr = "文件路径或地址错误: %s"
	catchErr  = " config.toml 不存在对应配置: %s"
)

var ef formatter.EpubFormat

func init() {
	ef.Init()
	rootCmd.PersistentFlags().StringVarP(&ef.Inner.Cover, coverCmd, "c", ef.Inner.Cover, "封面路径")
	rootCmd.PersistentFlags().StringVarP(&ef.Out, outCmd, "o", "", "输出文件名")
	rootCmd.PersistentFlags().StringVarP(&ef.Inner.ContentLogo, subCmd, "s", ef.Inner.ContentLogo, "每章标题logo")
	rootCmd.PersistentFlags().StringVarP(&ef.Inner.IntroImg, descImgCmd, "e", ef.Inner.IntroImg, "内容介绍logo")
	rootCmd.PersistentFlags().StringVarP(&ef.Inner.VolImg, volCmd, "b", ef.Inner.VolImg, "卷logo")
	rootCmd.PersistentFlags().StringVarP(&ef.Author, authorCmd, "a", "Unknown", "作者")
	rootCmd.PersistentFlags().StringVarP(&ef.Lang, langCmd, "l", "zh-Hans", "默认中文zh-Hans,英文 en")
	rootCmd.PersistentFlags().StringVarP(&ef.BookConf.Url, urlCmd, "i", "", "下载书籍介绍页(包含图片与简介页面) url")
	rootCmd.PersistentFlags().BoolVarP(&ef.BookConf.IsDesc, descCmd, "d", true, "是否包含制作说明,默认包含,使用 -d 来取消包含")
	rootCmd.PersistentFlags().StringVarP(&ef.BookConf.Path, pathCmd, "p", "", "转化txt路径")
	rootCmd.PersistentFlags().IntVarP(&ef.BookConf.Jump, jumpCmd, "j", 0, "跳过章节数")
	rootCmd.PersistentFlags().IntVarP(&ef.BookConf.Delay, delayCmd, "t", 0, "每章延迟毫秒数")
	rootCmd.PersistentFlags().StringVarP(&ef.BookConf.Catalog, catalogUrlCmd, "u", "", "章节爬取url 支持起点,番茄,七猫")
	rootCmd.PersistentFlags().StringVarP(&config.Cfg.From, configCmd, "f", "", "自定义 config.toml 路径(url或本地文件)")
}

var rootCmd = &cobra.Command{
	Use:   "freb",
	Short: "freb用于下载小说并转换至EPub",
	Run: func(cmd *cobra.Command, args []string) {
		path := ef.BookConf.Path
		if len(args) > 0 {
			path = args[0]
		}
		err := config.InitConfig()
		if err != nil {
			stdout.Errln(fmt.Errorf("配置错误: %v", err))
			return
		}
		bookCatch, err := CheckFlag(cmd, path)
		if err != nil {
			stdout.Errln(err)
			return
		}
		var source source.Source
		if len(ef.BookConf.Url) > 0 {
			source = &sources.UrlSource{}
		}
		if len(path) > 0 {
			source = &sources.TxtSource{}
			ef.Name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			ef.Name, ef.Author = utils.GetBookInfo(ef.Name)
		}

		err = InitAssets(ef, bookCatch)
		if err != nil {
			stdout.Errln(err)
			return
		}
		reg.InitTitleReg(bookCatch.Title.Filter)
		reg.InitContentReg(bookCatch.Content.Filter)

		err = source.GetBook(&ef, bookCatch)
		if err != nil {
			stdout.Errln(err)
			return
		}
		_ = os.RemoveAll(config.Cfg.TmpDir)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		stdout.Errln(err)
		os.Exit(1)
	}
}

func CheckFlag(cmd *cobra.Command, cmdPath string) (bookCatch *models.BookCatch, err error) {
	if cmd.PersistentFlags().Changed(descCmd) {
		ef.BookConf.IsDesc = false
	}
	if len(cmdPath) == 0 && len(ef.BookConf.Url) == 0 {
		return nil, errors.New(fromErr)
	}
	if len(cmdPath) > 0 {
		if !utils.IsFileInWorkDir(cmdPath) {
			return nil, errors.New(fmt.Sprintf(sourceErr, cmdPath))
		}
	}

	if !cmd.PersistentFlags().Changed(delayCmd) || ef.BookConf.Delay < 0 {
		ef.BookConf.Delay = config.Cfg.DelayTime
	}
	if len(ef.BookConf.Url) > 0 {
		for domain, catch := range config.Cfg.BookCatch {
			if strings.Contains(ef.BookConf.Url, domain) {
				bookCatch = catch
				bookCatch.Domain = domain
				return
			}
		}
		return nil, fmt.Errorf(catchErr, ef.BookConf.Url)
	}
	if len(config.Cfg.From) != 0 {
		if !reg.CheckUrl(config.Cfg.From) {
			return nil, errors.New(fmt.Sprintf(sourceErr, config.Cfg.From))
		}
	}
	return
}

func InitAssets(ef formatter.EpubFormat, bookCatch *models.BookCatch) (err error) {
	// image
	ef.Inner.Cover, err = utils.SetImage(ef.Cover, config.Cfg.TmpDir, ef.Inner.Cover, config.Cfg.From, func() *http.Request {
		var req *http.Request
		if len(ef.BookConf.Url) > 0 {
			var url, id string
			id, err = reg.MatchString(bookCatch.ID, ef.BookConf.Url)
			if err != nil {
				return nil
			}
			if bookCatch.Cover.NeedDivide {
				url = utils.DivideThousandURL(bookCatch.Cover.Url, id)
				for header, value := range bookCatch.Cover.Header {
					req.Header.Set(header, value)
				}
			} else {
				url = fmt.Sprintf(bookCatch.Cover.Url, id)
			}
			req = utils.GetWithUserAgent(url)
			return req
		}
		return nil
	})
	if err != nil {
		return
	}
	ef.Inner.IntroImg, err = utils.SetImage(ef.Inner.IntroImg, config.Cfg.TmpDir, ef.Inner.IntroImg, config.Cfg.From, nil)
	if err != nil {
		return
	}
	ef.Inner.ContentLogo, err = utils.SetImage(ef.Inner.ContentLogo, config.Cfg.TmpDir, ef.Inner.ContentLogo, config.Cfg.From, nil)
	if err != nil {
		return
	}
	ef.Inner.VolImg, err = utils.SetImage(ef.Inner.VolImg, config.Cfg.TmpDir, ef.Inner.VolImg, config.Cfg.From, nil)
	if err != nil {
		return
	}
	// font
	ef.AssetsPath.Font, err = utils.LocalOrDownload("assets/fonts/font.ttf", config.Cfg.From, config.Cfg.TmpDir)
	if err != nil {
		return
	}
	// MetaINF
	ef.AssetsPath.MetaInf, err = utils.LocalOrDownload("assets/META-INF/com.apple.ibooks.display-options.xml", config.Cfg.From, config.Cfg.TmpDir)
	if err != nil {
		return
	}
	// css
	ef.AssetsPath.MainCss, err = utils.LocalOrDownload("assets/styles/main.css", config.Cfg.From, config.Cfg.TmpDir)
	if err != nil {
		return
	}
	ef.AssetsPath.FontCss, err = utils.LocalOrDownload("assets/styles/fonts.css", config.Cfg.From, config.Cfg.TmpDir)
	if err != nil {
		return
	}
	ef.AssetsPath.CoverCss, err = utils.LocalOrDownload("assets/styles/cover.css", config.Cfg.From, config.Cfg.TmpDir)
	if err != nil {
		return
	}
	ef.AssetsPath.InstructionCss, err = utils.LocalOrDownload("assets/styles/instruction.css", config.Cfg.From, config.Cfg.TmpDir)
	if err != nil {
		return
	}
	return
}
