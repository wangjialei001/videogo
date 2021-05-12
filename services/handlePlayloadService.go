package services

import (
	"fmt"
	"regexp"
	"strings"
)

//type HandlePlayloadService struct{}

func test() string {
	fmt.Println("测试")
	return "ok"
}

//func (service *HandlePlayloadService) GetUrlParams(body string) []string {
func GetUrlParams(body string) []string {
	result := []string{}
	//解释正则表达式
	reg := regexp.MustCompile(`GET /hls/(?s:(.*?)) HTTP/1.1`)
	if reg == nil {
		fmt.Println("MustCompile err")
		return result
	}
	//提取关键信息
	items := reg.FindAllStringSubmatch(body, -1)
	//过滤<></>
	for _, text := range items {
		fmt.Println("text[1] = ", text[1])
		result = append(result, text[1])
	}
	return result
}
func GetUrlParamM3u8(body string) []string {
	result := []string{}
	//解释正则表达式
	reg := regexp.MustCompile(`GET /hls/(?s:(.*?)).m3u8 HTTP/1.1`)
	if reg == nil {
		fmt.Println("MustCompile err")
		return result
	}
	//提取关键信息
	items := reg.FindAllStringSubmatch(body, -1)
	//过滤<></>
	for _, text := range items {
		//fmt.Println("text[1] = ", text[1])
		var sParams = splitParam(text[1])
		for _, sParam := range sParams {
			result = append(result, sParam)
		}
	}
	return result
}
func GetUrlParamTs(body string) []string {
	result := []string{}
	//解释正则表达式
	reg := regexp.MustCompile(`GET /hls/(?s:(.*?)).ts HTTP/1.1`)
	if reg == nil {
		fmt.Println("MustCompile err")
		return result
	}
	//提取关键信息
	items := reg.FindAllStringSubmatch(body, -1)
	//过滤<></>
	for _, text := range items {
		//fmt.Println("text[1] = ", text[1])
		var sParams = splitParam(text[1])
		for _, sParam := range sParams {
			result = append(result, sParam)
		}
	}
	return result
}
func splitParam(parm string) []string {
	result := strings.Split(parm, "/video")
	// reg := regexp.MustCompile(`^[0-9]+(/video[0-9]+)*$`)
	// if reg == nil {
	// 	fmt.Println("splitParam MustCompile err")
	// 	return result
	// }
	// items := reg.FindAllStringSubmatch(parm, -1)
	// for _, item := range items {
	// 	result = append(result, item[1])
	// }
	return result
}
