package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
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
	"github.com/google/uuid"
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

		//API Routing configuration
		router.POST("/api/v1/texts", TextsController)

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

		//Gin port
		router.Run(":8080")
	}()

	EdgePath := "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe"
	cmd := exec.Command(EdgePath, "--app=http:127.0.0.1:8080/static/index.html")
	cmd.Start()
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)

	select {
	case <-chSignal:
		cmd.Process.Kill()
	}
}

func TextsController(c *gin.Context) {
	var json struct {
		Raw string `json:"raw"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		exe, err := os.Executable() // 获取当前执行文件的路径
		if err != nil {
			log.Fatal(err)
		}
		dir := filepath.Dir(exe) // 获取当前执行文件的目录
		if err != nil {
			log.Fatal(err)
		}
		filename := uuid.New().String()          // 生成一个文件名
		uploads := filepath.Join(dir, "uploads") // 拼接 uploads 的绝对路径
		err = os.MkdirAll(uploads, os.ModePerm)  // 创建 uploads 目录
		if err != nil {
			log.Fatal(err)
		}
		fullpath := path.Join("uploads", filename+".txt")                            // 拼接文件的绝对路径（不含 exe 所在目录）
		err = ioutil.WriteFile(filepath.Join(dir, fullpath), []byte(json.Raw), 0644) // 将 json.Raw 写入文件
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"url": "/" + fullpath}) // 返回文件的绝对路径（不含 exe 所在目录）
	}

}

func FilesController(c *gin.Context) {
	file, err := c.FormFile("raw")
	if err != nil {
		log.Fatal(err)
	}
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(exe)
	if err != nil {
		log.Fatal(err)
	}
	filename := uuid.New().String()
	uploads := filepath.Join(dir, "uploads")
	err = os.MkdirAll(uploads, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	fullpath := path.Join("uploads", filename+filepath.Ext(file.Filename))
	fileErr := c.SaveUploadedFile(file, filepath.Join(dir, fullpath))
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	c.JSON(http.StatusOK, gin.H{"url": "/" + fullpath})
}
