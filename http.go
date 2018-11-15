package main

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/text/encoding/simplifiedchinese"
)

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
