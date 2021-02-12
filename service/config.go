package service

import (
	"encoding/json"
	"file-server-go/model"
	"github.com/donething/utils-go/dofile"
	"log"
	"os"
)

// 配置文件名
var ConfigFile = ""

// 初始化配置
var Conf = model.Config{RootDir: ".", Auth: ""}

// 读取配置文件
func init() {
	if len(os.Args) <= 1 {
		log.Printf("需要指定配置文件的路径参数")
		os.Exit(0)
	}

	ConfigFile = os.Args[1]
	log.Printf("配置文件的路径：%s\n", ConfigFile)

	exists, err := dofile.Exists(ConfigFile)
	if !exists {
		// 触发配置文件不存在错误时，创建它
		saveConfig()
		log.Printf("请填写配置文件后，重新运行本程序\n")
		os.Exit(0)
	}

	// 读取配置
	data, err := dofile.Read(ConfigFile)
	if err != nil {
		log.Printf("读取配置文件(%s)出错：%s\n", ConfigFile, err)
		return
	}

	errParse := json.Unmarshal(data, &Conf)
	if errParse != nil {
		log.Printf("解析配置文件(%s)错误：%v\n", ConfigFile, errParse)
		return
	}
	log.Printf("已读取配置：%+v\n", Conf)
}

func saveConfig() bool {
	data, err := json.Marshal(Conf)
	if err != nil {
		log.Printf("将结构体配置数据(%+v)转为json格式数据失败：%v\n", Conf, err)
		return false
	}
	_, err = dofile.Write(data, ConfigFile, dofile.OTrunc, 0644)
	if err != nil {
		log.Printf("保存配置到文件失败：%v\n", err)
		return false
	}
	log.Printf("配置保存完成\n")
	return true
}
