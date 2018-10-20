package main

import (
	"bufio"
	"fmt"
	"github.com/Chyroc/download"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "163 open course download"
	app.Action = func(c *cli.Context) error {
		if len(c.Args()) == 0 {
			return cli.ShowAppHelp(c)
		}
		return run(c.Args().Get(0))
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(url string) error {
	courses, err := getCourseList(url)
	if err != nil {
		return err
	} else if len(courses) == 0 {
		return fmt.Errorf("未找到课程")
	}

	for index, course := range courses {
		fmt.Printf("%d\t%s\n", index, course.title)
	}
	fmt.Println()
	fmt.Printf("输入要下载的文件序号（%d-%d）下载该文件\n输入1,10,21下载序号为1和10和21的文件\n输入all选在下载所有文件\n", 0, len(courses)-1)

	reader := bufio.NewReader(os.Stdin)

	for {
		bs, _, err := reader.ReadLine()
		if err != nil {
			return err
		}

		var fileIndex []int
		var line = strings.TrimSpace(string(bs))
		if line == "all" {
			for i := range courses {
				fileIndex = append(fileIndex, i)
			}
		} else {
			gets := strings.Split(line, ",")
			for _, v := range gets {
				if v == "" {
					continue
				}
				i, err := strconv.Atoi(strings.TrimSpace(v))
				if err != nil {
					return fmt.Errorf("不合法的序号：%s", line)
				}
				fileIndex = append(fileIndex, i)
			}
		}

		for _, index := range fileIndex {
			fmt.Printf("下载：%s ...\n", courses[index].title)
			url, err := getDownloadURL(courses[index].url)
			if err != nil {
				return err
			}

			ext := filepath.Ext(url)

			if err = download.Download(url, "/tmp/"+courses[index].title+ext, 20); err != nil {
				return err
			}
		}
	}

	return nil
}
