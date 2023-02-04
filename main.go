package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func get(url string) []string {
	var images []string

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	rxp := regexp.MustCompile(`<img[^>]+src=["']([^"']+)["']`)

	submatch := rxp.FindAllStringSubmatch(string(b), -1)
	for _, item := range submatch {
		images = append(images, item[1])
	}

	return images
}

func download(target string) {
	if _, err := os.Stat("images"); os.IsNotExist(err) {
		os.Mkdir("images", os.ModeDir)
	}
	res, err := http.Get(target)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	targetSplit := strings.Split(target, "/")
	name := targetSplit[len(targetSplit)-1]

	fp := filepath.Join("images", name)
	fmt.Println(fp)

	f, err := os.Create(fp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if _, err := f.Write(b); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func main() {
	args := os.Args[1:]

	if len(args) > 1 || len(args) < 1 {
		fmt.Println("invalid parameters")
		return
	}

	host := args[0]

	if host[len(host)-1] != '/' {
		host = host + "/"
	}

	imgs := get(host)

	for _, img := range imgs {
		download(img)
	}
}
