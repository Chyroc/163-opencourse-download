package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Chyroc/download"
	"github.com/urfave/cli"
)

var dir string

func main() {
	app := cli.NewApp()
	app.Name = "163 open course download"
	app.Action = func(c *cli.Context) error {
		if len(c.Args()) == 0 {
			return cli.ShowAppHelp(c)
		}
		return run(c.Args().Get(0))
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "dir",
			Usage:       "视频存放位置",
			Destination: &dir,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(url string) error {
	if dir == "" {
		return fmt.Errorf("请指定 -dir")
	}

	if err := ensureDirExist(dir); err != nil {
		return err
	}

	courses, err := getCourseList(url)
	if err != nil {
		return err
	} else if len(courses) == 0 {
		return fmt.Errorf("未找到课程")
	}

	for index, course := range courses {
		fmt.Printf("%d\t%s\n", index, course.title)
	}
	fmt.Printf(`
输入要下载的文件序号（%d-%d）下载该文件
输入"1,10,21"下载序号为1和10和21的文件
输入"5-8"下载序号为5、6、7、8的文件
输入"all"选在下载所有文件

`, 0, len(courses)-1)

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
			for _, comma := range strings.Split(line, ",") {
				if strings.Index(comma, "-") != -1 {
					barList := strings.Split(comma, "-")
					if len(barList) != 2 {
						return fmt.Errorf("不合法的序号：%s", line)
					}
					start := strings.TrimSpace(barList[0])
					end := strings.TrimSpace(barList[1])
					if start == "" || end == "" {
						return fmt.Errorf("不合法的序号：%s", line)
					}
					a, err := strconv.Atoi(start)
					if err != nil {
						return fmt.Errorf("不合法的序号：%s", line)
					}
					b, err := strconv.Atoi(end)
					if err != nil {
						return fmt.Errorf("不合法的序号：%s", line)
					}

					for i := a; i <= b; i++ {
						fileIndex = append(fileIndex, i)
					}
				} else {
					v := strings.TrimSpace(comma)
					if v == "" {
						continue
					}
					i, err := strconv.Atoi(v)
					if err != nil {
						return fmt.Errorf("不合法的序号：%s", line)
					}
					fileIndex = append(fileIndex, i)
				}
			}
		}

		for _, index := range fileIndex {
			fmt.Printf("下载：%s ...\n", courses[index].title)
			url, err := getDownloadURL(courses[index].url)
			if err != nil {
				return err
			}

			savefile := filepath.Join(dir, courses[index].title+filepath.Ext(url))
			if err = download.Download(url, savefile, 20); err != nil {
				return err
			}
		}
	}

	return nil
}
