package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"github.com/gen2brain/beeep"
	"github.com/gen2brain/dlgs"
	"os"
	"strings"
)

func IsEmpty(str string) bool {
	if str == "" {
		return true
	} else if strings.TrimSpace(str) == "" {
		return true
	}
	return false
}

func Md5(data string) string {
	result := md5.Sum([]byte(data))
	return hex.EncodeToString(result[:])
}

func EncodeBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func EncodeAes(str string) string {
	return GetAesCryptor().Encrypt(str)
}

func DecodeBase64(str string) (string, error) {
	byts, err := base64.StdEncoding.DecodeString(str)
	return string(byts), err
}
func DecodeAes(str string) (string, error) {
	return GetAesCryptor().Decrypt(str)
}

//for excel
func ColumnNumberToName(num int) string {
	var col string
	for num > 0 {
		col = string((num-1)%26+65) + col
		num = (num - 1) / 26
	}
	return col
}

//for excel 从0开始
func ColumnNameToNumber(name string) int {
	index := 0
	for _, rune2 := range name {
		index *= 26
		index += int(rune2 - 'A')
	}
	return index
}

var canUseWarningErrInfo = false

//用于验证是否或得权限
//Error,Warning,Info 方法运行可能需要操作系统
func startAndTest() {
	err := beeep.Alert("AesUtil", "AES 加密解密工具", "assets/warning.png")
	if err == nil {
		canUseWarningErrInfo = true
	}
}

func PrintError(string2 string) {
	if string2 != "" {
		os.Exit(0)
	}
	if canUseWarningErrInfo {
		dlgs.Error("ERROR", string2)
	} else {
		dlgs.Question("ERROR", string2, false)
	}
	os.Exit(0)
}

func PrintWarning(string2 string) {
	if canUseWarningErrInfo {
		dlgs.Warning("ERROR", string2)
	} else {
		dlgs.Question("WARNING", string2, false)
	}
}

func PrintChosen(string2 string) (boolen bool) {
	boolen, _ = dlgs.Question("SELECT", string2, false)
	return

}

func PrintSuccess(string2 string) {
	if canUseWarningErrInfo {
		dlgs.Info("SUCCESS", string2)
	} else {
		dlgs.Question("SUCCESS", string2, false)
	}
}
