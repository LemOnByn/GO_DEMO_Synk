package main

import (
	"io/fs"
	"log"
	"net/http"
	"strings"

	// "log"
	// "net/url"
	// "fmt"
	"embed"
	"os"
	"os/exec"
	"os/signal"

	// "syscall"
	"github.com/gin-gonic/gin"
	// "github.com/zserge/lorca"
)

//将目录的内容全部打包进exe文件
//go:embed frontend/dist/*
var FS embed.FS

func main() {
	//gin 协程
	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()

		// router.GET("/", func(c *gin.Context) { //Context是gin的结构体 封装了Request 和 Writer(response)
		// 	c.Writer.Write([]byte("hello world"))
		// })

		//将静态文件变成变量
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		//静态文件置于static这个url下然后用http去读取解析静态文件
		router.StaticFS("/static", http.FS(staticFiles))
		//静态资源访问失败处理
		router.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/static/") {
				reader, err := staticFiles.Open("index.html")
				if err != nil {
					log.Fatal(err)
				}
				defer reader.Close()
				stat, err := reader.Stat()
				if err != nil {
					log.Fatal(err)
				}
				c.DataFromReader(http.StatusOK, stat.Size(), "text/html;charset=utf-8", reader, nil)
			} else {
				c.Status(http.StatusNotFound)
			}
		})

		router.Run(":8080")
	}()

	EdgePath := "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe"
	cmd := exec.Command(EdgePath, "--app=http:127.0.0.1:8080/static/index.html")
	cmd.Start()
git
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)

	select {
	case <-chSignal:
		cmd.Process.Kill()
	}
}
