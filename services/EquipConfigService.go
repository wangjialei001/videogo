package services

import (
	"encoding/json"
	"fmt"
	"gopacket/model"
	"io"
	"os"
	"path"
	"runtime"
)

func ReadJsons() []model.EquipVideoConfigModel {
	const dataFile = "config/equipVideo.json"
	_, filename, _, _ := runtime.Caller(1)
	//fmt.Println(filename)
	datapath := path.Join(path.Dir(filename), dataFile)
	fmt.Println(datapath)
	f, error := os.Open(datapath)
	if error != nil {
		panic(error)
	}
	defer f.Close()
	r, error := io.ReadAll(f)
	if error != nil {
		panic(error)
	}
	//fmt.Printf("json数据：%s", string(r))
	equipVideoConfig := make([]model.EquipVideoConfigModel, 0)
	err := json.Unmarshal([]byte(r), &equipVideoConfig)
	if err != nil {
		panic(err)
	}
	return equipVideoConfig
}
