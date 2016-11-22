package byteoperator

//处理 []int32 的工具类
import (
	"bytes"
)

func ToBytes(data []int32) []byte {
	buf := bytes.NewBuffer([]byte{})
	length := len(data)
	buf.Write(Int32ToBytes(int32(length)))

	for _, v := range data {
		buf.Write(Int32ToBytes(v))
	}

	return buf.Bytes()
}

func FromBytes(data []byte) []int32 {
	start := 0
	_, datalen := Bytes2Int32(data, start)

	start += 4
	result := make([]int32, 0)
	var i int32 = 0
	for i = 0; i < datalen; i++ {
		_, v32 := Bytes2Int32(data, start)
		start += 4
		result = append(result, v32)
	}

	return result
}
