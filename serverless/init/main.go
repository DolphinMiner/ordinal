package main

import (
	"embed"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
)

func main() {
	// lambda.Start(Handler)
	Handler()
}

//go:embed punks
var pngFS embed.FS

func Handler() {
	for pngIndex := 0; pngIndex < 1000; pngIndex++ {
		fileName := getPngName(pngIndex)
		fmt.Printf("current filename is %s", fileName)
		bytes, err := pngFS.ReadFile(fileName)
		if err != nil {
			log.Fatal("fail to open local png resource!", err)
		}
		// 转十六进制
		hexStr := hex.EncodeToString(bytes)
		// 初始化db
		fmt.Printf("get hex png resource is %s", hexStr)
	}
}

func getPngName(index int) string {
	suffix := strconv.Itoa(index)
	append_zore := 4 - len(suffix)
	for i := 0; i < append_zore; i++ {
		suffix = "0" + suffix
	}
	return "punks/" + suffix + ".png"
}
