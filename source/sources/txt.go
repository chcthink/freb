package sources

import (
	"bufio"
	"bytes"
	"fmt"
	"freb/formatter"
	"freb/models"
	"freb/utils/reg"
	"freb/utils/stdout"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	titleMax     = 32
	unknownTitle = "未知章节"
)

type TxtSource struct {
}

func getBuffer(filename string) *bufio.Reader {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		stdout.Errfln("读取文件出错: %s", err.Error())
		os.Exit(1)
	}
	temBuf := bufio.NewReader(f)
	bs, _ := temBuf.Peek(1024)
	encodig, encodename, _ := charset.DetermineEncoding(bs, "text/plain")
	if encodename != "utf-8" {
		f.Seek(0, 0)
		bs, err := io.ReadAll(f)
		if err != nil {
			stdout.Errfln("读取文件出错: %s", err.Error())
			os.Exit(1)
		}
		var buf bytes.Buffer
		decoder := encodig.NewDecoder()
		if encodename == "windows-1252" {
			decoder = simplifiedchinese.GB18030.NewDecoder()
		}
		bs, _, _ = transform.Bytes(decoder, bs)
		buf.Write(bs)
		return bufio.NewReader(&buf)
	} else {
		f.Seek(0, 0)
		buf := bufio.NewReader(f)
		return buf
	}
}

func (t *TxtSource) GetBook(ef *formatter.EpubFormat, catch *models.BookCatch) error {
	var contentList []models.Section

	stdout.Fmtln("正在读取txt文件...")
	start := time.Now()
	buf := getBuffer(ef.BookConf.Path)
	var title string
	content := &bytes.Buffer{}
	tmp := ""
	isIntro := false
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if line = strings.TrimSpace(line); line != "" && !checkComment(line) {
					ef.GenLine2Buffer(line, content)
				}
				contentList = append(contentList, models.Section{
					Title:   title,
					Content: content.String(),
				})
				content.Reset()
				break
			}
			return fmt.Errorf("读取文件出错: %s", err)
		}
		line = strings.TrimSpace(line)
		// 空行直接跳过
		if len(line) == 0 {
			continue
		}
		// 跳过注释行
		if checkComment(line) {
			if len(contentList) == 0 {
				content.Reset()
			}
			continue
		}
		if len(contentList) == 0 {
			if isAuthor, author := reg.GetAuthor(line); isAuthor {
				ef.Author = author
				if tmp != "" {
					ef.Name = tmp
				}
				continue
			}
			tmp = line
			intro := ""
			if !isIntro {
				if isIntro, intro = reg.GetIntro(line); isIntro {
					if intro != "" {
						ef.Intro = intro
					}
					continue
				}
			}
			if !(reg.CheckTitle(line) || reg.CheckVol(line)) && isIntro {
				ef.Intro += line
			} else {
				isIntro = false
				content.Reset()
			}
		}

		// 处理标题
		if utf8.RuneCountInString(line) <= titleMax &&
			(reg.CheckTitle(line) || reg.CheckVol(line) || reg.CheckEnd(line)) {
			if title == "" {
				title = unknownTitle
			}
			if content.Len() > 0 || title != unknownTitle {
				contentList = append(contentList, models.Section{
					Title:   title,
					Content: content.String(),
				})
			}
			title = line
			content.Reset()
			continue
		}
		if line == "（全书完）" {
			contentList = append(contentList, models.Section{
				Title:   title,
				Content: content.String(),
			})
			content.Reset()
			break
		}
		ef.GenLine2Buffer(line, content)
	}
	// 没识别到章节又没识别到 EOF 时，把所有的内容写到最后一章
	if content.Len() != 0 {
		if title == "" {
			title = "章节正文"
		}
		contentList = append(contentList, models.Section{
			Title:   title,
			Content: content.String(),
		})
	}
	ef.Sections = contentList
	err := ef.InitBook()
	if err != nil {
		return err
	}
	var volPath string
	for i := range ef.Sections {
		volPath, err = ef.GenBookContent(i, volPath)
	}
	err = ef.Build()
	if err != nil {
		return err
	}
	end := time.Now().Sub(start)
	stdout.Successfln("\n已生成书籍,使用时长: %s", end)
	return nil
}

// checkComment 判断是否为备注,形如: =======  ////// ***** -----
func checkComment(content string) bool {
	if strings.ReplaceAll(content, "=", "") == "" {
		return true
	}
	if strings.ReplaceAll(content, "*", "") == "" {
		return true
	}
	if strings.ReplaceAll(content, "-", "") == "" {
		return true
	}
	if strings.ReplaceAll(content, "/", "") == "" {
		return true
	}
	return false
}
