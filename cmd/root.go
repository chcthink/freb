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
	"freb/utils/htmlx"
	"freb/utils/reg"
	"freb/utils/stdout"
	"github.com/spf13/cobra"
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
	searchCmd     = "search"
)

const (
	fromErr   = "数据来源异常(-p -i 皆为空)"
	sourceErr = "文件路径或地址错误: %s"
	catchErr  = " config.toml 不存在对应配置: %s"
)

var searchValue string

var ef formatter.EpubFormat

func init() {
	ef.Init()
	rootCmd.PersistentFlags().StringVarP(&ef.Images.Cover, coverCmd, "c", ef.Images.Cover, "封面路径")
	rootCmd.PersistentFlags().StringVarP(&ef.Out, outCmd, "o", "", "输出文件名")
	rootCmd.PersistentFlags().StringVarP(&ef.Images.ContentLogo, subCmd, "s", ef.Images.ContentLogo, "每章标题logo")
	rootCmd.PersistentFlags().StringVarP(&ef.Images.IntroImg, descImgCmd, "e", ef.Images.IntroImg, "内容介绍logo")
	rootCmd.PersistentFlags().StringVarP(&ef.Images.VolImg, volCmd, "b", ef.Images.VolImg, "卷logo")
	rootCmd.PersistentFlags().StringVarP(&ef.Author, authorCmd, "a", "Unknown", "作者")
	rootCmd.PersistentFlags().StringVarP(&ef.Lang, langCmd, "l", "zh-Hans", "默认中文zh-Hans,英文 en")
	rootCmd.PersistentFlags().StringVarP(&ef.BookConf.Url, urlCmd, "i", "", "下载书籍介绍页(包含图片与简介页面) url")
	rootCmd.PersistentFlags().BoolVarP(&ef.BookConf.IsDesc, descCmd, "d", true, "是否包含制作说明,默认包含,使用 -d 来取消包含")
	rootCmd.PersistentFlags().StringVarP(&ef.BookConf.Path, pathCmd, "p", "", "转化txt路径")
	rootCmd.PersistentFlags().IntVarP(&ef.BookConf.Jump, jumpCmd, "j", 0, "跳过章节数")
	rootCmd.PersistentFlags().IntVarP(&ef.BookConf.Delay, delayCmd, "t", 0, "每章延迟毫秒数")
	rootCmd.PersistentFlags().StringVarP(&ef.BookConf.Catalog, catalogUrlCmd, "u", "", "章节爬取url 支持起点,番茄,七猫")
	rootCmd.PersistentFlags().StringVarP(&config.Cfg.From, configCmd, "f", "", "自定义 config.toml 路径(url或本地文件)")
	rootCmd.PersistentFlags().StringVarP(&searchValue, searchCmd, "g", "", "在已配置域名下搜索书名")
}

var rootCmd = &cobra.Command{
	Use:   "freb",
	Short: "freb用于下载小说并转换至EPub",
	Run: func(cmd *cobra.Command, args []string) {
		path := ef.BookConf.Path
		if len(args) > 0 {
			path = args[0]
		}
		err := InitAssets(ef)
		if err != nil {
			stdout.Errln(err)
			return
		}
		if searchValue != "" {
			ef.BookConf.Url, err = source.Search(searchValue, config.Cfg.BookCatch)
			if err != nil {
				return
			}
		}
		bookCatch, err := checkFlag(cmd, path)
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
			ef.Name, ef.Author = reg.GetBookInfo(ef.Name)
		}

		err = source.GetBook(&ef, bookCatch)
		if err != nil {
			stdout.Errln(err)
			return
		}

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		stdout.Errln(err)
		os.Exit(1)
	}
}

func checkFlag(cmd *cobra.Command, cmdPath string) (bookCatch *models.BookCatch, err error) {
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
	if len(config.Cfg.From) != 0 {
		if !reg.CheckUrl(config.Cfg.From) {
			return nil, errors.New(fmt.Sprintf(sourceErr, config.Cfg.From))
		}
	}
	if len(ef.BookConf.Url) > 0 {
		return CheckCatchUrl(ef.BookConf.Url)
	}
	return
}

func CheckCatchUrl(url string) (bookCatch *models.BookCatch, err error) {
	if !reg.CheckUrl(url) {
		return nil, errors.New(fmt.Sprintf(sourceErr, url))
	}
	for domain, catch := range config.Cfg.BookCatch {
		if strings.Contains(url, domain) {
			bookCatch = catch
			bookCatch.Domain = domain
			reg.InitCustomFilterReg(bookCatch)
			return
		}
	}
	return nil, fmt.Errorf(catchErr, url)
}

func InitAssets(ef formatter.EpubFormat) (err error) {
	ef.Images.IntroImg, err = htmlx.DownloadWithReq(ef.Images.IntroImg, config.Cfg.TmpDir, config.Cfg.From, nil)
	if err != nil {
		return
	}
	ef.Images.ContentLogo, err = htmlx.DownloadWithReq(ef.Images.ContentLogo, config.Cfg.TmpDir, config.Cfg.From, nil)
	if err != nil {
		return
	}
	ef.Images.VolImg, err = htmlx.DownloadWithReq(ef.Images.VolImg, config.Cfg.TmpDir, config.Cfg.From, nil)
	if err != nil {
		return
	}
	// font
	ef.Assets.Font, err = htmlx.LocalOrDownload("assets/fonts/font.ttf", config.Cfg.From, config.Cfg.TmpDir)
	if err != nil {
		return
	}
	// MetaINF
	ef.Assets.MetaInf, err = htmlx.LocalOrDownload("assets/META-INF/com.apple.ibooks.display-options.xml", config.Cfg.From, config.Cfg.TmpDir)
	if err != nil {
		return
	}
	// css
	ef.Assets.MainCss, err = htmlx.LocalOrDownload("assets/styles/main.css", config.Cfg.From, config.Cfg.TmpDir)
	if err != nil {
		return
	}
	ef.Assets.FontCss, err = htmlx.LocalOrDownload("assets/styles/fonts.css", config.Cfg.From, config.Cfg.TmpDir)
	if err != nil {
		return
	}
	ef.Assets.CoverCss, err = htmlx.LocalOrDownload("assets/styles/cover.css", config.Cfg.From, config.Cfg.TmpDir)
	if err != nil {
		return
	}
	ef.Assets.InstructionCss, err = htmlx.LocalOrDownload("assets/styles/instruction.css", config.Cfg.From, config.Cfg.TmpDir)
	if err != nil {
		return
	}
	return
}
