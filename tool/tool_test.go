package tool

import "testing"

func TestMD5Str(t *testing.T) {
	t.Log(MD5Str(""))
	t.Log(MD5Str("123"))
}

func TestBytesHumanReadable(t *testing.T) {
	t.Log(BytesHumanReadable(1023))
	t.Log(BytesHumanReadable(1025))
	t.Log(BytesHumanReadable(1024032))
	t.Log(BytesHumanReadable(102403200000))
}
