package fileutils

import (
	"encoding/json"
	"log"
	"os"
)

func OpenJsonFileWithParam(fileDir string, result interface{}) interface{} {
	r, err := os.Open(fileDir)
	if err != nil {
		log.Fatalln("open json file error:" + err.Error())
	}
	decoder := json.NewDecoder(r)

	err = decoder.Decode(&result)
	if err != nil {
		log.Fatalln("JSON Decode error:" + err.Error())
	}
	return result
}
