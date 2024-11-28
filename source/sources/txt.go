package sources

import (
	"bufio"
	"bytes"
	"fmt"
	"freb/formatter"
	"freb/models"
	"freb/utils"
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

func readBuffer(filename string) *bufio.Reader {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		fmt.Println("读取文件出错: ", err.Error())
		os.Exit(1)
	}
	temBuf := bufio.NewReader(f)
	bs, _ := temBuf.Peek(1024)
	encodig, encodename, _ := charset.DetermineEncoding(bs, "text/plain")
	if encodename != "utf-8" {
		f.Seek(0, 0)
		bs, err := io.ReadAll(f)
		if err != nil {
			fmt.Println("读取文件出错: ", err.Error())
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

func (t *TxtSource) GetBook(book *models.Book) error {
	var contentList []models.Chapter
	var a []string
	var ef formatter.EpubFormat
	ef.Book = book
	err := ef.InitBook()
	if err != nil {
		return err
	}
	fmt.Println("正在读取txt文件...")
	start := time.Now()
	buf := readBuffer(book.Path)
	var title string
	content := &bytes.Buffer{}
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if line = strings.TrimSpace(line); line != "" {
					ef.GenLine2Buffer(line, content)
				}
				contentList = append(contentList, models.Chapter{
					Title:   title,
					Content: content.String(),
				})
				content.Reset()
				break
			}
			return fmt.Errorf("读取文件出错: %w", err)
		}
		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, "<", "&lt;")
		line = strings.ReplaceAll(line, ">", "&gt;")
		// 空行直接跳过
		if len(line) == 0 {
			continue
		}
		// 处理标题
		if utf8.RuneCountInString(line) <= titleMax &&
			(utils.CheckTitle(line) || utils.CheckVol(line)) {
			if title == "" {
				title = unknownTitle
			}
			if content.Len() > 0 || title != unknownTitle {
				fmt.Println(content)
				contentList = append(contentList, models.Chapter{
					Title:   title,
					Content: content.String(),
				})
				a = append(a, content.String())
			}
			title = line
			content.Reset()
			continue
		}
		ef.GenLine2Buffer(line, content)
	}
	// 没识别到章节又没识别到 EOF 时，把所有的内容写到最后一章
	if content.Len() != 0 {
		if title == "" {
			title = "章节正文"
		}
		contentList = append(contentList, models.Chapter{
			Title:   title,
			Content: content.String(),
		})
	}
	// var sectionList []models.Chapter
	book.Chapters = contentList
	var volPath string
	for i := range book.Chapters {
		volPath, err = ef.GenBookContent(i, volPath)
		// if utils.CheckVol(section.Title) {
		// 	if volumeSection != nil {
		// 		sectionList = append(sectionList, *volumeSection)
		// 		volumeSection = nil
		// 	}
		// 	volumeSection = &section
		//
		// } else {
		// 	if volumeSection == nil {
		// 		sectionList = append(sectionList, section)
		// 	} else {
		// 		volumeSection.Sections = append(volumeSection.Sections, section)
		// 	}
		// }
	}
	// 如果有最后一卷,添加到章节列表
	// if volumeSection != nil {
	// 	sectionList = append(sectionList, *volumeSection)
	// 	volumeSection = nil
	// }
	err = ef.Build()
	if err != nil {
		return err
	}
	end := time.Now().Sub(start)
	fmt.Println("\n已生成书籍,使用时长: ", end)
	// fmt.Println("匹配章节:", sectionCount(sectionList))
	return nil
}

func sectionCount(sections []models.Chapter) int {
	var count int
	for _, section := range sections {
		count += 1 + len(section.Sections)
	}
	return count
}
