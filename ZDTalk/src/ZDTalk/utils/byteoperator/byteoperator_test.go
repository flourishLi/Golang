package byteoperator

import (
	"fmt"
	//	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	//	data := Int64ToBytes(1234567890123456789)

	fmt.Println(Bytes2HexString(Int64ToBytes(1234567890123456789)))
}

func Test_Float32ToBytes(t *testing.T) {
	data := Int32ToBytes(1000000000)
	data2 := Float32ToBytes(0.6)
	fmt.Println(data)
	fmt.Println(data2)
	_, v1 := Bytes2Int32(data, 0)
	err, v2 := Bytes2Float32(data, 0)
	fmt.Println(v1)
	fmt.Println(v2, err)
	fmt.Println(1.11 > 1.1)
}
