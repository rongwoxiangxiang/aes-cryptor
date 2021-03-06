package main

import (
	"github.com/gen2brain/dlgs"
	"os"
	"strings"
)

func chooseFileSingle() {
	defer recoverTop()
	boolen, _ := dlgs.Question("请选择文件", "请选择需要解密的文件，当前仅支持 txt/xls/xlsx 格式的文件,请选择文件", false)
	if !boolen {
		checkClose()
	}
	decodeFilePhoneOrIdCard()
}

func decodeFilePhoneOrIdCard() {
	filePath := getFile()
	_, err := os.Stat(filePath)
	if err == nil {
		mode, _, err := dlgs.List("解密内容", "请选择解密内容:", []string{"身份证", "手机号"})
		if err != nil {

		}
		if mode == "身份证" {
			needBase64Bool = true
		}
		if strings.Contains(filePath, ".xls") {
			decodeExcel(filePath)
		} else if strings.Contains(filePath, ".txt") {
			decodeTxt(filePath)
		} else {
			PrintError("当前仅支持 txt/xls/xlsx 格式的文件")
		}
	}
}
