package service

import (
	"file-server-go/logger"
	. "file-server-go/model"
	"file-server-go/tool"
	"github.com/donething/utils-go/dotext"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ListFiles 获取指定路径下的文件列表
func ListFiles(path string) JResult {
	filesList, err := ioutil.ReadDir(path)
	if err != nil {
		logger.Error.Printf("获取路径(%s)下的文件列表出错：%v\n", path, err)
		return JResult{Success: false, Code: 21, Msg: "获取文件列表出错：" + err.Error()}
	}
	// 将返回的文件列表
	files := make([]FileDesp, 0, len(filesList)) // 将返回的文件描述的切片

	// 按修改时间排序（最近下载的排到前面）
	sort.Slice(filesList, func(i, j int) bool {
		return filesList[i].ModTime().After(filesList[j].ModTime())
	})

	// 获取信息
	for i := 0; i < len(filesList); i++ {
		f := FileDesp{
			Name:  filesList[i].Name(),
			Last:  dotext.FormatDate(filesList[i].ModTime(), dotext.TimeFormatDefault),
			Size:  tool.BytesHumanReadable(filesList[i].Size()),
			IsDir: filesList[i].IsDir(),
		}
		files = append(files, f)
	}
	logger.Info.Printf("返回路径(%s)下的文件列表\n", path)
	return JResult{Success: true, Code: 10, Msg: "文件列表", Data: files}
}

// DelFile 删除文件
func DelFile(path string) JResult {
	if err := os.RemoveAll(path); err != nil {
		logger.Error.Printf("删除文件(%s)失败：%v\n", path, err)
		return JResult{Success: false, Code: 22, Msg: "删除文件失败：" + err.Error()}
	}
	logger.Info.Printf("删除文件(%s)成功\n", path)
	return JResult{Success: true, Code: 10, Msg: "删除文件成功", Data: path}
}

// CheckPath 检查路径是否合法或越界，不允许访问所设根目录更前的位置
func CheckPath(path string) (string, bool) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		logger.Error.Printf("访问的路径有误(%s)，校正为所设的根目录(%s)\n", path, Conf.Root)
		return Conf.Root, false
	}
	rootAbsPath, _ := filepath.Abs(Conf.Root)
	if !strings.HasPrefix(absPath, rootAbsPath) {
		logger.Warn.Printf("访问的路径非法(%s)，校正为所设的根目录(%s)\n", path, Conf.Root)
		return Conf.Root, false
	}
	return absPath, true
}
