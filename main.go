package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gen2brain/dlgs"
	"strings"
)

var (
	operator string
)

func main() {
	startAndTest()
	loginStatus := false
	for trys := 0; trys < 5; trys++ {
		if login() { //登录失败重试
			loginStatus = true
			break
		}
		if !PrintChosen("输入错误，请重新输入") {
			PrintError("错误次数过多")
		}
	}
	if loginStatus {
		//chooseFileSingle 只提供文件解密功能  chooseMode()提供字符串文件加密解密功能
		chooseFileSingle()
		//chooseMode()
	}
}

//Y2Njcy4xMTEyMjIy
func login() bool {
	keys, boolen, err := dlgs.Password("Password", "请输入账户密钥")
	if err != nil {
		loggerError(err)
	}
	if !boolen {
		checkClose()
	}
	loginDatas, err := DecodeBase64(keys)
	if err != nil || !strings.Contains(loginDatas, ".") {
		return false
	}
	loginData := strings.Split(loginDatas, ".")
	name := loginData[0]
	pass := loginData[1]
	operator = name
	if validLogin(name, pass) {
		return true
	}
	operator = ""
	return false
}

func chooseMode() {
	defer recoverTop()
	mode, _, err := dlgs.List("解密数据为文件/字符串", "请选择解密内容类型：文件、字符串", []string{"FILE", "STRINGS"})
	if err != nil {
		loggerError(err)
	}
	if mode == "FILE" {
		processFile()
	} else {
		processStrings()
	}
}

func processStrings() {
	var printStr, errString string
	var data map[string]string
	mode, _, err := dlgs.List("模式", "请选择模式:", []string{"decode-aes", "decode-aes+base64", "encode-ase", "encode-base64+ase"})
	if err != nil {
		loggerError(err)
	}
	for {
		data, errString = doProcessStings(mode)
		if errString != "" {
			PrintWarning(errString)
		}
		if data != nil || len(data) < 1 {
			for key, val := range data {
				if !IsEmpty(key) {
					printStr += key + " ==>>> " + val + "\n"
				}
			}
			if printStr != "" {
				PrintSuccess(printStr)
			}
		}
		if !needContinue() {
			break
		}
		printStr = ""
	}
}

func doProcessStings(mode string) (data map[string]string, errString string) {
	var (
		strs   string
		err    error
		boolen bool
	)
	for {
		if strings.Contains(mode, "encode") {
			strs, boolen, err = dlgs.Entry("请输入", "请输入需要加密的字符串，多个字符串以逗号分隔:", "")
		} else {
			strs, boolen, err = dlgs.Entry("请输入", "请输入需要解密的字符串，多个字符串以逗号分隔", "")
		}
		if !boolen {
			checkClose()
		}
		if err != nil {
			loggerError(err)
		}
		strs = strings.TrimSpace(strs)
		if !IsEmpty(strs) {
			break
		}
	}

	strArr := strings.Split(strs, ",")
	switch mode {
	case "decode-aes":
		data, errString = decodeAesString(strArr)
	case "decode-aes+base64":
		data, errString = decodeAesBase64String(strArr)
	case "encode-ase":
		data = encodeAESString(strArr)
	case "encode-base64+ase":
		data = encodeBase64AesString(strArr)
	}
	return
}

func processFile() {
	PrintMessage("当前仅支持 txt/xls/xlsx 格式的文件,请选择文件")
	mode := getAESMode()
	if mode == "AES_ENCODE" {
		//TODO
		PrintError("暂不支持 AES 文件加密")
	} else {
		decodeFile()
	}
}

func validLogin(name, pass string) bool {
	user, err := new(User).GetByName(name)
	if err != nil {
		return false
	}
	signCalc := base64.StdEncoding.EncodeToString([]byte(pass + user.Salt))
	if IsEmpty(user.Pass) || Md5(signCalc) != user.Pass {
		logger("PASSWORD_ERR", "LOGIN")
		return false
	}
	logger("SUCCESS", "LOGIN")
	return true
}

func recoverTop() {
	if r := recover(); r != nil {
		PrintError(fmt.Sprintf("AES ERROR %v", r))
	}
}
