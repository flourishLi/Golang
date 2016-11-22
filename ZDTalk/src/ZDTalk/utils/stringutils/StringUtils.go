package stringutils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func String2int(value string) int {
	v, _ := strconv.Atoi(value)
	return v
}
func String2int32(value string) int32 {
	if StringIsEmpty(value) {
		return -1
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return -1
	}
	return int32(v)

}
func String2int64(value string) int64 {
	v64, _ := strconv.ParseInt(value, 10, 64)
	return v64
}

//判断字符串是否为空
func StringIsEmpty(value string) bool {
	value = strings.TrimSpace(value)
	return len(value) == 0 || strings.EqualFold(value, "")
}

//判断字节数组是否为空
func byteIsEmpty(value []byte) bool {
	return len(value) == 0 || value == nil
}

/*获取当前文件执行的路径*/
func GetCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	splitstring := strings.Split(path, "\\")
	size := len(splitstring)
	splitstring = strings.Split(path, splitstring[size-1])
	ret := strings.Replace(splitstring[0], "\\", "/", -1)
	fmt.Print(ret)
	return ret
}

// 获取当前文件执行的文件名
func GetCurFilename() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	//	splitstring := strings.Split(path, "\\")
	splitstring := strings.Split(path, string(filepath.Separator))
	size := len(splitstring)
	filenamestrings := strings.Split(splitstring[size-1], ".")
	filenameOnly := filenamestrings[0]

	return filenameOnly
}

func GetFileName(filePath string) string {
	splitstring := strings.Split(filePath, "/")
	size := len(splitstring)
	//	filenamestrings := strings.Split(splitstring[size-1], "/")
	//	filenameOnly := filenamestrings[size-1]
	filenameOnly := splitstring[size-1]

	return filenameOnly
}

//截取字符串
func SubString(value string, start, length int) string {
	rs := []rune(value)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

//切片转换成字符串，以oper分隔
func SliceToString(slices interface{}, oper string) (string, error) {
	var buf bytes.Buffer
	first := true
	switch value := slices.(type) {
	case []int32:
		for _, v := range value {
			if first {
				first = false
			} else {
				buf.WriteString(oper)
			}
			buf.WriteString(fmt.Sprintf("%v", v))
		}
	case []string:
		for _, v := range value {
			if first {
				first = false
			} else {
				buf.WriteString(oper)
			}
			buf.WriteString(fmt.Sprintf("%v", v))
		}
	case []bool:
		for _, v := range value {
			if first {
				first = false
			} else {
				buf.WriteString(oper)
			}
			buf.WriteString(fmt.Sprintf("%v", v))
		}
	default:
		return "", errors.New("Slice type errors,only the basic types")
	}
	return buf.String(), nil
}

func GetStringIds(source string) []int32 {
	var data []int32

	v := strings.Split(source, ",")
	fmt.Println("id source:" + source)
	for _, i := range v {
		v, e := strconv.Atoi(i)
		if e == nil {
			fmt.Println("id:", v)
			data = append(data, int32(v))
		}
	}

	return data
}

func GetIdsToString(ids []int64) string {
	data := ""
	for _, i := range ids {
		if data != "" {
			data += ","
		}

		data += strconv.FormatInt(i, 10)

	}

	return data
}

//获取基本类型数据格式化成string的值
func GetFormatString(in interface{}) string {
	var s string
	switch in.(type) {
	case int:
		s = fmt.Sprintf("%d", in)
	case int16:
		s = fmt.Sprintf("%d", in)
	case int32:
		s = fmt.Sprintf("%d", in)
	case int64:
		s = fmt.Sprintf("%d", in)
	case float32:
		s = fmt.Sprintf("%f", in)
	case float64:
		s = fmt.Sprintf("%f", in)
	default:
		return s
	}

	return s
}
