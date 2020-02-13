package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/gen2brain/beeep"
	"github.com/gen2brain/dlgs"
	"os"
	"reflect"
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

//for excel A -> 0 B->1
func ColumnNumberToName(num int) string {
	var col string
	for num > 0 {
		col = string((num)%26+65) + col
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
	err := beeep.Alert("AesUtil", "AES UTIL", "assets/warning.png")
	if err == nil {
		canUseWarningErrInfo = true
	}
}

//数组平分
func SplitArray(arr interface{}, num int) (segmens []int) {
	lens := reflect.ValueOf(arr).Cap()
	if lens < num {
		return append(segmens, 0)
	}
	var step = lens / num
	for end := 1; end <= lens; end += step {
		segmens = append(segmens, end-1)
	}
	return
}

//按步数分割数组
func SplitArrayStep(arr interface{}, step int) (segmens []int) {
	lens := reflect.ValueOf(arr).Cap()
	if lens < step {
		return append(segmens, 0)
	}
	for end := 1; end <= lens; end += step {
		segmens = append(segmens, end-1)
	}
	return
}

func PrintError(string2 string) {
	fmt.Println(string2)
	if canUseWarningErrInfo {
		dlgs.Error("ERROR", string2)
	} else {
		dlgs.Question("ERROR", string2, false)
	}
	logger(string2, "ERROR")
	os.Exit(0)
}

func PrintWarning(string2 string) {
	fmt.Println(string2)
	if canUseWarningErrInfo {
		beeep.Alert("Warning", string2, "assets/warning.png")
	} else {
		dlgs.Question("WARNING", string2, false)
	}
}

func PrintChosen(string2 string) (boolen bool) {
	boolen, _ = dlgs.Question("SELECT", string2, false)
	return
}

func PrintSuccess(string2 string) {
	fmt.Println(string2)
	if canUseWarningErrInfo {
		beeep.Alert("SUCCESS", string2, "assets/warning.png")
	} else {
		dlgs.Question("SUCCESS", string2, false)
	}
}

func PrintMessage(string2 string) {
	fmt.Println(string2)
	if canUseWarningErrInfo {
		beeep.Alert("MESSAGE", string2, "assets/warning.png")
	} else {
		dlgs.Question("MESSAGE", string2, false)
	}
}

func Notify(string2 string) {
	if canUseWarningErrInfo {
		beeep.Alert("Notify", string2, "assets/warning.png")
	} else {
		dlgs.Question("Notify", string2, false)
	}
}

func logger(content, operation string) {
	fmt.Println(content, operation)
	log := new(Log)
	log.Content = content
	log.Operator = operator
	log.Operation = operation
	_, err := log.Insert(nil)
	if err != nil {
		PrintError("未知错误，" + err.Error())
	}
}

func loggerError(err error) {
	fmt.Println(err.Error())
	log := new(Log)
	log.Content = err.Error()
	log.Operator = operator
	log.Operation = "ERROR"
	_, errr := log.Insert(nil)
	if errr != nil {
		PrintError("未知错误，" + errr.Error())
	}
	os.Exit(0)
}

func getFile() string {
	file, _, err := dlgs.File("请选择文件", "", false)
	if err != nil {
		loggerError(err)
	}
	return file
}

func getAESMode() string {
	mode, _, err := dlgs.List("模式", "请选择模式:", []string{"AES_ENCODE", "AES_DECODE"})
	if err != nil {
		loggerError(err)
	}
	return mode
}

func needBase64DecodeEncode() bool {
	needBase64Decode, err := dlgs.Question("二次解密/加密", "是否需要base64二次解密/加密", false)
	if err != nil {
		loggerError(err)
	}
	return needBase64Decode
}

func needContinue() bool {
	needCont, err := dlgs.Question("是否继续操作", "是否继续操作", false)
	if err != nil {
		loggerError(err)
	}
	return needCont
}
func checkClose() {
	needClose, err := dlgs.Question("是否关闭", "是否关闭AES工具", false)
	if err != nil {
		loggerError(err)
	}
	if needClose {
		os.Exit(0)
	}
}
