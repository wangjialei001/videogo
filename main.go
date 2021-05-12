package main

import (
	"flag"
	"fmt"
	"gopacket/services"
	"gopacket/utils"
	"strconv"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {

	version := pcap.Version()
	fmt.Println(version)
	fmt.Println("packet start...")

	equipVideoConfigs := services.ReadJsons()
	if equipVideoConfigs == nil {
		fmt.Println("equipVideoConfig is nil...")
		return
	}
	fmt.Printf("配置文件：%v", equipVideoConfigs)
	videoProcMaps := make(map[uint64]int)
	//设备配置信息
	equipVideoMap := make(map[uint64]time.Time)
	go func() {
		for {
			//times++
			//fmt.Println("tick", times)
			time.Sleep(time.Minute * time.Duration(1))
			now := time.Now()
			for equip, expireTime := range equipVideoMap {
				//fmt.Printf("设备 %d： 预计失效时间：%v", equip, expireTime)
				//fmt.Print(equip)
				//fmt.Print(expireTime)
				fmt.Printf("设备 %d 预计失效时间：%s\n", equip, expireTime.Format("2006-01-02 15:04:05"))
				if now.After(expireTime) {
					//删除失效设备
					delete(equipVideoMap, equip)
					//关闭设备视频
					delete(videoProcMaps, equip)
					utils.Get("http://localhost/api/video/closevideo?equipId=" + strconv.Itoa(int(equip)))
					fmt.Printf("设备 %d 预计失效时间：%s，已经关闭\n", equip, expireTime.Format("2006-01-02 15:04:05"))
				}
			}
		}
	}()
	deviceName := ""
	flag.StringVar(&deviceName, "device", "ens33", "输入网卡")
	var portInt int
	flag.IntVar(&portInt, "port", 80, "输入监控端口")
	flag.Parse()
	snapLen := int32(65535)
	filter := getFilter(uint16(portInt))
	fmt.Println("device:%v, snapLen:%v, port:%s\n", deviceName, snapLen, portInt)
	fmt.Println("filter:", filter)
	//打开网络接口，抓取在线数据
	handle, err := pcap.OpenLive(deviceName, snapLen, true, pcap.BlockForever)
	if err != nil {
		fmt.Printf("pcap open live failed: %v", err)
		return
	}
	// 设置过滤器
	if err := handle.SetBPFFilter(filter); err != nil {
		fmt.Printf("set bpf filter failed: %v", err)
		return
	}
	defer handle.Close()
	// 抓包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetSource.NoCopy = true
	for packet := range packetSource.Packets() {
		if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			fmt.Println("unexpected packet")
			continue
		}

		//fmt.Printf("packet:%v\n", packet)

		// tcp 层
		//tcp, _ := packet.TransportLayer().(*layers.TCP)

		//fmt.Printf("tcp:%v\n", tcp)
		// tcp payload，也即是tcp传输的数据
		//fmt.Printf("tcp payload:%v\n", tcp.Payload)
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		equipIdStr := ""
		if ipLayer != nil {
			tcpLayer := packet.Layer(layers.LayerTypeTCP)
			if tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)
				if len(tcp.Payload) > 0 {
					//handlePlayloadService := HandlePlayloadService
					// urlParams := handlePlayloadService.GetUrlParams(string(tcp.Payload))
					// for _, text := range urlParams {
					// 	fmt.Println("text[1] = ", text)
					// }

					urlParams := services.GetUrlParamM3u8(string(tcp.Payload))
					equipIdStr = ""
					if len(urlParams) >= 2 {
						equipIdStr = urlParams[0]
						//fmt.Printf("开始观看视频，视频目录：%s，文件名称：%s.m3u8\n，", urlParams[0], urlParams[1])
					}
					urlParamTss := services.GetUrlParamTs(string(tcp.Payload))
					if len(urlParamTss) >= 2 {
						equipIdStr = urlParamTss[0]
						//fmt.Printf("视频ts文件，视频目录：%s，文件名称：%s.ts\n", urlParamTss[0], urlParamTss[1])
					}

					if len(equipIdStr) > 0 {
						//fmt.Printf("当前设备Id:%s\n", equipIdStr)
						intEquipId, intEquipIdError := strconv.Atoi(equipIdStr)
						if intEquipIdError != nil {
							continue
						}
						int64EquipId := uint64(intEquipId)
						afterM, afterMError := time.ParseDuration("2m") //设置为2分钟失效
						if afterMError != nil {
							continue
						}
						equipVideoMap[int64EquipId] = time.Now().Add(afterM)
						//判断是否已经播放
						if _, videoProcMapOk := videoProcMaps[int64EquipId]; !videoProcMapOk {
							for _, equipVideoConfig := range equipVideoConfigs { //查看设备配置信息
								if equipVideoConfig.EquipId-int64EquipId == 0 {
									fmt.Printf("设备Id:%d，等待开启视频\n", int64EquipId)
									//videoProcMaps[int64EquipId] = services.StartVideo(equipVideoConfig)
									//表示已经开启视频播放
									videoProcMaps[int64EquipId] = 1
									//fmt.Printf("设备Id:%d，进程Id：%d", int64EquipId, videoProcMaps[int64EquipId])
									utils.Get("http://localhost/api/video/showvideo?equipId=" + strconv.Itoa(int(int64EquipId))) //目前依赖.net core 程序执行shell命令
									fmt.Printf("设备Id:%d，已经开启视频\n", int64EquipId)
									break
								}
							}
						}
					}

					// ip, _ := ipLayer.(*layers.IPv4)
					// fmt.Printf("源ip %s:%s\n",
					// 	ip.SrcIP, tcp.SrcPort)
					// fmt.Printf("目的ip %s:%s\n",
					// 	ip.DstIP, tcp.DstPort)
					// fmt.Printf("参数3 %s\n",
					// 	string(tcp.Payload))
				}
			} else if errLayer := packet.ErrorLayer(); errLayer != nil {
				fmt.Printf("tcp.err: %v", errLayer)
			}
		} else if errLayer := packet.ErrorLayer(); errLayer != nil {
			fmt.Printf("ip.err:%v", errLayer)
		}
	}

}
func getFilter(port uint16) string {
	filter := fmt.Sprintf("tcp and ((src port %v) or (dst port %v))", port, port)
	return filter
}
