package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type course struct {
	url   string
	title string
}

var downloadURL = regexp.MustCompile(`appsrc : '(.*?)',`)

func getDownloadURL(url string) (string, error) {
	html, err := httpGet(url)
	if err != nil {
		return "", err
	}

	match := downloadURL.FindStringSubmatch(html)
	if len(match) < 2 {
		return "", fmt.Errorf("不合法的url：%s", url)
	}

	return strings.Replace(match[1], "-list.m3u8", ".flv", -1), err
}

func getCourseList(url string) ([]course, error) {
	if strings.Contains(url, "open.163.com/movie") {
		return getCourseListOfMovie(url)
	} else if strings.Contains(url, "open.163.com/special") {
		return getCourseListOfSpecial(url)
	}

	return nil, fmt.Errorf("不合法的链接，请访问：https://open.163.com/")
}

func getCourseListOfMovie(url string) ([]course, error) {
	html, err := httpGet(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var courses []course
	doc.Find("#j-playlist-container > div > div").Each(func(i int, s *goquery.Selection) {
		if s.HasClass("j-hoverdown") {
			// have url
			url := s.Find("a").AttrOr("href", "")
			title := s.Find("p > span").Text() + s.Find("p > a").Text()
			courses = append(courses, course{url: url, title: title})
			return
		}

		courses = append(courses, course{url: url, title: s.Find("p.f-thide").Text()})
	})

	return courses, nil
}

func getCourseListOfSpecial(url string) ([]course, error) {
	html, err := httpGet(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var courses []course
	doc.Find("#list2 > tbody > tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		title := strings.TrimSpace(s.Find("td.u-ctitle").Text())
		title = strings.Replace(title, " ", "", -1)
		title = strings.Replace(title, "\n", "", -1)
		url := s.Find("td.u-ctitle > a").AttrOr("href", "")

		courses = append(courses, course{url: url, title: title})
	})

	return courses, nil
}

func httpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	bs, err = gbkToUtf8(bs)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func gbkToUtf8(s []byte) ([]byte, error) {
	return simplifiedchinese.GBK.NewDecoder().Bytes(s)
}

func lastFilename(filename string) string {
	i := strings.LastIndex(filename, "/")
	if i == -1 {
		return filename
	}
	return filename[i:]
}
