package timeutils

import "time"

func Unix64TimeToUnix32Time(data int64) int32 {

	// value := data >> 32
	return (int32)(data)
}
func GetUnix13NowTime() int64 {
	unixtime := time.Now().UTC().UnixNano()
	return unixtime / 1000000
}

func GetTimeStamp() int32 {
	return int32(GetUnix13NowTime() / 1000)
}

func GetUnixTimeStamp() int64 { //秒
	unixtime := time.Now().UTC().UnixNano()
	return unixtime / 1000000000
}
