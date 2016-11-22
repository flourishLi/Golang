package byteoperator

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/***
byte转int8
**/
func Bytes2Int8(data []byte, start int) (error, int8) {

	err := checkDataLength(data, start, 1)

	if err != nil {
		return err, 0
	}

	tmp := data[start : start+1]

	b_buf := bytes.NewBuffer(tmp)

	var result int8
	binary.Read(b_buf, binary.BigEndian, &result)

	return nil, result
}

/***
byte转int16
**/
func Bytes2Int16(data []byte, start int) (error, int16) {

	err := checkDataLength(data, start, 2)

	if err != nil {
		return err, 0
	}

	tmp := data[start : start+2]

	b_buf := bytes.NewBuffer(tmp)

	var result int16
	binary.Read(b_buf, binary.BigEndian, &result)

	return nil, result
}

/***
byte转int32
**/
func Bytes2Int32(data []byte, start int) (error, int32) {

	err := checkDataLength(data, start, 4)

	if err != nil {
		return err, 0
	}

	tmp := data[start : start+4]

	b_buf := bytes.NewBuffer(tmp)

	var result int32
	binary.Read(b_buf, binary.BigEndian, &result)

	return nil, result
}

/***
byte转int64
**/
func Bytes2Int64(data []byte, start int) (error, int64) {
	err := checkDataLength(data, start, 8)

	if err != nil {
		return err, 0
	}

	tmp := data[start : start+8]

	b_buf := bytes.NewBuffer(tmp)

	var result int64
	binary.Read(b_buf, binary.BigEndian, &result)

	return nil, result
}

/***
byte转float32
**/
func Bytes2Float32(data []byte, start int) (error, float32) {
	err := checkDataLength(data, start, 4)

	if err != nil {
		return err, 0
	}

	tmp := data[start : start+4]

	b_buf := bytes.NewBuffer(tmp)

	var result float32
	binary.Read(b_buf, binary.BigEndian, &result)

	return nil, result
}

func Bytes2String(data []byte, start int, length int) (error, string) {

	err := checkDataLength(data, start, length)

	if err != nil {
		return err, ""
	}

	tmp := data[start : start+length]
	index := bytes.IndexByte(tmp, 0)
	if index > 0 {
		tmp = tmp[0:index]
	}
	//fmt.Println("tmp---->", tmp)
	result := string(tmp)
	result = strings.TrimSpace(result)
	return nil, result
}

func Bytes2HexString(data []byte) string {

	buf := bytes.NewBuffer([]byte{})

	count := 1
	for _, d := range data {
		s := fmt.Sprintf("%0x", d)

		if len(s) == 1 {
			s = "0" + s
		}

		buf.WriteString(s)

		if count%4 == 0 {
			buf.WriteString(" ")
		}
		count++
	}

	return string(buf.Bytes())
}

/***
int8转byte
**/
func Int8ToBytes(data int8) []byte {
	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, data)
	return b_buf.Bytes()
}

/***
int16转byte
**/
func Int16ToBytes(data int16) []byte {
	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, data)

	return b_buf.Bytes()
}

/***
int32转byte
**/
func Int32ToBytes(data int32) []byte {
	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, data)
	return b_buf.Bytes()
}

/***
int64转byte
**/
func Int64ToBytes(data int64) []byte {
	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, data)
	return b_buf.Bytes()
}

/***
float32转byte
**/
func Float32ToBytes(data float32) []byte {
	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, data)
	return b_buf.Bytes()
}

/**
*字符串转换为字节数组
 */
func String2BytesWithIndex(data string, start int, length int) []byte {
	if len(data) == 0 {
		return []byte{}
	}
	len2 := len(data) - start

	if len2 > length {
		len2 = length
	}

	tmp2 := []byte(data[start:len2])
	result := tmp2[0:]
	if len2 < length {
		len3 := length - len(result)
		tmp3 := make([]byte, len3)
		result = append(result, tmp3...)
	}
	return result
}

/**
*字符串转换为字节数组
 */
func String2Bytes(data string) []byte {

	result := []byte(data[0:])
	//	tmp3 := Int32ToBytes(int32(len2))
	len2 := len(data)
	tmp3 := Int16ToBytes(int16(len2))

	result = append(tmp3, result...)
	return result
}

func BytesToStringIncludeDataLength(data []byte, start int) (error, int16, string) {
	//	_, length := Bytes2Int32(data, start)
	//	_, result := Bytes2String(data, start+4, int(length))
	_, length := Bytes2Int16(data, start)

	if length == 0 {
		return nil, length, ""
	}

	if length <= 0 || length >= int16(0x8000-1) {
		return errors.New("length error:" + strconv.Itoa(int(length))), 0, ""
	}
	_, result := Bytes2String(data, start+2, int(length))

	return nil, length, result
}

func BytesToMsgBytesIncludeDataLength(data []byte, start int) (error, int16, []byte) {
	_, length := Bytes2Int16(data, start)
	if length == 0 {
		return nil, length, data
	}

	if length <= 0 || length >= int16(0x8000-1) {
		return errors.New("length error:" + strconv.Itoa(int(length))), 0, []byte{}
	}
	result := data[start+2 : start+2+int(length)]
	return nil, length, result
}

func MsgBytesToBytes(data []byte) []byte {
	result := data
	len2 := len(data)
	tmp3 := Int16ToBytes(int16(len2))
	result = append(tmp3, result...)
	return result
}

//检查数据长度
func checkDataLength(data []byte, start int, length int) error {
	if len(data)-start < length || start < 0 {
		return errors.New("数据长度不对")
	}
	return nil
}
