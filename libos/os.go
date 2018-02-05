package OS

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"strconv"
	"strings"
)

const (
	PATH_DELIMITER = "/"
)

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//生成Guid字串
func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

//检测路径结束符;
func CheckPathDelimiter(path string) string {
	if 0 == len(path) {
		return ""
	} else if (len(path) - 1) == strings.LastIndex(path, PATH_DELIMITER) {
		return path
	}
	return (path + PATH_DELIMITER)
}

//字符串拼接;
func CreateText(args ...interface{}) (string, error) {
	var buf bytes.Buffer
	for _, val := range args {
		switch val.(type) {
		case int:
			buf.WriteString(strconv.Itoa(val.(int)))
		case int32:
			buf.WriteString(strconv.FormatInt(int64(val.(int32)), 10))
		case int64:
			buf.WriteString(strconv.FormatInt(val.(int64), 10))
		case string:
			buf.WriteString(val.(string))
		default:
			return "", errors.New("Unknown Type")
		}
	}
	return buf.String(), nil
}
