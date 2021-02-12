package service

import (
	"encoding/json"
	"file-server-go/model"
	"github.com/donething/utils-go/dofile"
	"log"
	"os"
)

const (
	// 配置文件名
	ConfigFile = `fileserver.conf`
)

// 初始化配置
var Conf = model.Config{RootDir: ".", Auth: ""}

func init() {
	path := ConfigFile
	exists, err := dofile.Exists(path)
	if !exists {
		// 触发配置文件不存在错误时，创建它
		saveConfig()
		log.Printf("请填写配置文件后，重新运行本程序\n")
		os.Exit(0)
	}

	data, err := dofile.Read(path)
	if err != nil {
		log.Printf("读取配置文件(%s)出错：%s\n", path, err)
		return
	}

	errParse := json.Unmarshal(data, &Conf)
	if errParse != nil {
		log.Printf("解析配置文件(%s)错误：%v\n", path, errParse)
		return
	}
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
