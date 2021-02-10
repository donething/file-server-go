package tool

import (
	"crypto/md5"
	"fmt"
)

// 返回字符串的 md5 值
func MD5Str(str string) (m string) {
	m5 := md5.Sum([]byte(str))
	m = fmt.Sprintf("%x", m5)
	return
}

// 将字节转为可读的字符（如"102MB"）
// https://stackoverflow.com/a/30822306
func BytesHumanReadable(bytes int64) string {
	var kb int64 = 1024
	mb := 1024 * kb
	gb := 1024 * mb
	tb := gb * 1024
	if bytes < kb {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes >= kb && bytes < mb {
		return fmt.Sprintf("%d KB", bytes/kb)
	} else if bytes >= mb && bytes < gb {
		return fmt.Sprintf("%d MB", bytes/mb)
	} else if bytes >= gb && bytes < tb {
		return fmt.Sprintf("%d GB", bytes/gb)
	} else {
		return fmt.Sprintf("%d TB", bytes/tb)
	}
}
