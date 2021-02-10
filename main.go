package main

import (
	. "file-server-go/model"
	. "file-server-go/service"
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
	router.Use(CORS("http://127.0.0.1:10012"))

	// 使用 gzip 压缩，语句"gzip.Gzip(gzip.DefaultCompression)"不能放在Middleware()中，否则无效
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(TokenAuthMiddleware)

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusForbidden, "")
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
	// 当操作为浏览、下载时，不需验证
	if c.Query("op") == "list" || c.Query("op") == "down" {
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

// 准许指定域名的 CORS
func CORS(host string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", host)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
