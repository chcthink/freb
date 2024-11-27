package cmd

import (
	"fmt"
	"freb/config"
	"freb/models"
	"freb/source"
	"freb/utils"
	"github.com/spf13/cobra"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	formatCmd = "format"
	idCmd     = "id"
	coverCmd  = "cover"
	volCmd    = "vol"
	outCmd    = "out"
	subCmd    = "sub"
	authorCmd = "author"
	descCmd   = "desc"
	langCmd   = "lang"
	// fontCmd   = "font"
)

const (
	coverDefault    = "cover.png"
	subCoverDefault = "sub_cover.png"
	volDefault      = "vol.png"
)

const (
	flagErr    = "flag 值异常: %s"
	textUrlErr = "图片地址错误: %s"
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&novel.Format, formatCmd, "f", "epub", "转换至指定格式 默认 epub")
	rootCmd.PersistentFlags().StringVarP(&novel.Id, idCmd, "i", "", "下载书本id")
	rootCmd.PersistentFlags().StringVarP(&novel.Cover, coverCmd, "c", coverDefault, "封面路径")
	rootCmd.PersistentFlags().StringVarP(&novel.Out, outCmd, "o", "", "输出文件名")
	rootCmd.PersistentFlags().StringVarP(&novel.SubCover, subCmd, "s", subCoverDefault, "每章封面")
	rootCmd.PersistentFlags().StringVarP(&novel.Vol, volCmd, "b", volDefault, "卷图片")
	rootCmd.PersistentFlags().StringVarP(&novel.Author, authorCmd, "a", "Unknown", "作者")
	rootCmd.PersistentFlags().BoolVarP(&novel.Desc, descCmd, "d", true, "是否包含制作说明,默认包含,使用 -d 来取消包含")
	rootCmd.PersistentFlags().StringVarP(&novel.Lang, langCmd, "l", "zh-Hans", "默认中文zh-Hans,英文 en")
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
		var source source.UrlSource
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
	// format
	switch novel.Format {
	case "epub":
	case "txt":
	default:
		return fmt.Errorf(flagErr, formatCmd)
	}
	// id
	if novel.Id == "" && utils.NumReg(novel.Id) {
		return fmt.Errorf(flagErr, idCmd)
	}
	// cover subCover vol
	novel.Cover, err = setImage(novel.Cover, "cover", true, func() *http.Request {
		bookId, _ := strconv.Atoi(novel.Id)
		mid := strconv.FormatFloat(math.Floor(float64(bookId)/1000.0), 'f', 0, 64)

		url := strings.Join([]string{utils.Domain(), "fengmian", mid, novel.Id, novel.Id + "s.jpg"}, "/")
		req := utils.NewGet(url)
		req.Header.Set("Referer", "https://69shuba.cx/modules/article/search.php")
		return req
	})
	if err != nil {
		return err
	}
	novel.SubCover, err = setImage(novel.SubCover, "sub_cover", false, func() *http.Request {
		return utils.NewGet(novel.SubCover)
	})
	if err != nil {
		return err
	}
	novel.Vol, err = setImage(novel.Vol, "vol", false, func() *http.Request {
		return utils.NewGet(novel.Vol)
	})
	if err != nil {
		return err
	}
	return nil
}

func setImage(from, filename string, ifDefaultReq bool, handler func() *http.Request) (source string, err error) {
	if utils.IsImgFile(from) {
		if utils.IsFileExist(from) {
			return
		}
		if ifDefaultReq {
			source, err = utils.DownloadTmp(config.Cfg.TmpDir, filename, handler)
			if err != nil {
				utils.Warnf(textUrlErr, err.Error())
			}
		}
		return
	}
	if utils.CheckUrl(from) {
		source, err = utils.DownloadTmp(config.Cfg.TmpDir, "cover", func() *http.Request {
			return utils.NewGet(novel.Cover)
		})
		if err != nil {
			utils.Warnf(textUrlErr, err.Error())
		}
	}
	return
}
