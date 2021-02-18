package main

import (
	. "file-server-go/model"
	. "file-server-go/service"
	"fmt"
	"github.com/donething/utils-go/dofile"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	// 输出的格式
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile) // 打印 log 时显示时间戳
}

func main() {
	port := "20200"
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// 准许指定域名的 CORS
	// 域名为 nginx 中当前 web 网站的监听地址和端口
	//router.Use(CORS("http://127.0.0.1:10012"))
	router.Use(Cors())

	// 使用 gzip 压缩，语句"gzip.Gzip(gzip.DefaultCompression)"不能放在Middleware()中，否则无效
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(TokenAuthMiddleware)

	router.GET("/", func(c *gin.Context) {
		log.Printf("禁止访问\n")
		c.String(http.StatusForbidden, "禁止访问")
	})
	router.GET("/api/file", FileHandler)

	log.Printf("开始服务：http://127.0.0.1:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// /api/file?op=list&rel=.
func FileHandler(c *gin.Context) {
	operator := c.Query("op")
	rel := c.Query("rel")

	// 路径越界检查
	path := filepath.Join(Conf.RootDir, rel)
	log.Printf("收到的路径参数：%s\n", path)

	if path, allow := CheckPath(path); !allow {
		log.Printf("禁止访问更前面的路径(%s)\n", path)
		c.JSON(http.StatusOK, JResult{Success: false, Code: 20, Msg: "禁止访问更前面的路径"})
		return
	}
	// 获取指定路径的信息
	fi, err := os.Stat(path)
	if err != nil {
		log.Printf("获取文件(%s)信息时出错：%s\n", path, err)
		c.JSON(http.StatusOK, JResult{Success: false, Code: 21, Msg: "获取文件信息时：" + err.Error()})
	}

	// 判断执行的操作
	switch strings.ToLower(operator) {
	case "list":
		c.JSON(http.StatusOK, ListFiles(path))
		break
	case "del":
		c.JSON(http.StatusOK, DelFile(path))
		break
	case "down":
		if fi.IsDir() {
			log.Printf("指定的文件路径(%s)为目录，不支持下载\n", path)
			c.JSON(http.StatusOK, JResult{Success: false, Code: 22, Msg: "指定的文件路径为目录"})
			return
		}
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+fi.Name())
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Accept-Length", fmt.Sprintf("%d", fi.Size()))
		log.Printf("开始提供路径(%s)的下载\n", path)
		c.File(path)
		break
	case "md5":
		md5, err := dofile.Md5(path)
		if err != nil {
			log.Printf("计算文件(%s)的md5值出错：%s\n", path, err)
			c.JSON(http.StatusOK, JResult{Success: false, Code: 23, Msg: "计算文件md5值出错"})
			return
		}
		c.JSON(http.StatusOK, JResult{Success: true, Code: 10, Msg: "文件的MD5值", Data: md5})
		break
	default:
		log.Printf("未知的操作：%s: %s\n", operator, rel)
		c.JSON(http.StatusOK, JResult{Success: false, Code: 24, Msg: "未知的操作"})
	}
}

// 授权验证中间件
func TokenAuthMiddleware(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	token := Conf.Auth
	// 当操作不为删除文件时不需验证
	if c.Query("op") != "del" {
		c.Next()
		return
	}

	if auth != token {
		// c.AbortWithStatus(http.StatusUnauthorized)
		log.Printf("授权验证错误：'%s'\n", auth)
		c.String(http.StatusUnauthorized, "未授权的访问")
		c.Abort()
	}
}

// 准许 CORS
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// 允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}
