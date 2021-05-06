package service

import (
	"encoding/json"
	"file-server-go/logger"
	"file-server-go/model"
	"github.com/donething/utils-go/dofile"
	"os"
)

// ConfigPath 默认的配置文件
var ConfigPath = "file-server.json"

// Conf 初始化配置
var Conf = model.Config{Root: "", Auth: ""}

// 读取配置文件
func init() {
	if len(os.Args) >= 2 {
		ConfigPath = os.Args[1]
	}
	logger.Info.Printf("读取配置文件：%s\n", ConfigPath)

	exists, err := dofile.Exists(ConfigPath)
	if !exists {
		// 触发配置文件不存在错误时，创建它
		saveConfig()
		logger.Warn.Printf("已生成配置文件，请填写配置后，重新运行本程序\n")
		os.Exit(0)
	}

	// 读取配置
	data, err := dofile.Read(ConfigPath)
	if err != nil {
		logger.Error.Printf("读取配置文件(%s)出错：%s\n", ConfigPath, err)
		os.Exit(0)
	}

	err = json.Unmarshal(data, &Conf)
	if err != nil {
		logger.Error.Printf("解析配置文件(%s)错误：%v\n", ConfigPath, err)
		os.Exit(0)
	}
	logger.Info.Printf("已读取配置：%+v\n", Conf)

	// 判断配置是否有效
	if Conf.Root == "" {
		logger.Warn.Printf("配置中'root'路径为空，请先填写再运行本程序")
		os.Exit(0)
	}
}

func saveConfig() bool {
	data, err := json.Marshal(Conf)
	if err != nil {
		logger.Error.Printf("将结构体配置数据(%+v)转为json格式数据失败：%v\n", Conf, err)
		return false
	}
	_, err = dofile.Write(data, ConfigPath, dofile.OTrunc, 0644)
	if err != nil {
		logger.Error.Printf("保存配置到文件失败：%v\n", err)
		return false
	}
	return true
}
