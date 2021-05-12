package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

func Get(url string) string {
	client := http.Client{Timeout: 5 * time.Second}
	res, err := client.Get(url)
	if err != nil {
		fmt.Println(err.Error)
	}
	defer res.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := res.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		}
	}
	return result.String()
}
