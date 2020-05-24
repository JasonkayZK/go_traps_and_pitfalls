package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

/* defer in Copy file */
func main() {
	filename, n, _ := fetch(os.Args[1])
	fmt.Printf("download file: %s, size: %d bit\n", filename, n)
}

func fetch(url string) (filename string, n int64, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	// url末尾作为文件名
	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}
	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}

	// 文件大小
	n, err = io.Copy(f, resp.Body)

	// 不使用defer关闭文件
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}

	return local, n ,err
}
