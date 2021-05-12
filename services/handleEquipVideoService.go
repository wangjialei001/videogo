package services

import (
	"fmt"
	"gopacket/model"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func StopVideo(equipId uint64) {

}
func StartVideo(videoConfig model.EquipVideoConfigModel) int {
	var tsPath = "/opt/nginx-1.19.0/html/hls/" + strconv.Itoa(int((videoConfig.EquipId)))
	fmt.Println(tsPath)
	_, err := os.Stat(tsPath)
	if err == nil {
		return 0
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(tsPath, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error)
		}
	}
	var build strings.Builder
	build.WriteString("ffmpeg -rtsp_transport tcp -i  ")
	build.WriteString(`"rtsp://"`)
	build.WriteString(videoConfig.UserName)
	build.WriteString(":")
	build.WriteString(videoConfig.Pwd)
	build.WriteString("@")
	build.WriteString(videoConfig.Url)
	build.WriteString(":")
	build.WriteString(strconv.Itoa(int(videoConfig.Port)))
	build.WriteString("/Streaming/Channels/")
	build.WriteString(strconv.Itoa(int(videoConfig.Num)))
	build.WriteString("01\" -fflags flush_packets -max_delay 1 -flags -global_header -hls_init_time 0 -hls_time 1 -master_pl_publish_rate 5 -hls_list_size 50 -hls_flags delete_segments -vcodec copy -y /opt/nginx-1.19.0/html/hls/")
	build.WriteString(strconv.Itoa(int(videoConfig.EquipId)))
	build.WriteString("/video")
	build.WriteString(strconv.Itoa(int(videoConfig.EquipId)))
	build.WriteString(".m3u8")
	cmdStr := build.String()
	fmt.Printf("执行命令：%s", cmdStr)
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	errCmd := cmd.Start()
	if errCmd != nil {
		panic(errCmd)
	}
	fmt.Println("son ", cmd.Process.Pid)
	return cmd.Process.Pid
}
