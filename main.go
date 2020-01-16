package main

import (
	"fmt"
	"github.com/gen2brain/dlgs"
	"os"
	"strings"
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
			os.Exit(0)
		}

	}
	if loginStatus {
		chooseMode()
	}
}

//Y2Njcy4xMTEyMjIy
func login() bool {
	keys, boolen, err := dlgs.Password("Password", "请输入账户密钥")
	if err != nil {
		panic(err)
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
	if validLogin(name, pass) {
		return true
	}
	return false
}

func chooseMode() {
	defer recoverTop()
	mode, _, err := dlgs.List("解密数据为文件/字符串", "请选择解密内容类型：文件、字符串", []string{"FILE", "STRINGS"})
	if err != nil {
		panic(err)
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
		panic(err)
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
			panic(err)
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
	PrintWarning("当前仅支持 txt/xls/xlsx 格式的文件,请选择文件")
	mode := getAESMode()
	if mode == "AES_ENCODE" {
		//TODO
		PrintError("暂不支持")
	} else {
		decodeFile()
	}
}

func getFile() string {
	file, _, err := dlgs.File("请选择文件", "", false)
	if err != nil {
		panic(err)
	}
	return file
}

func getAESMode() string {
	mode, _, err := dlgs.List("模式", "请选择模式:", []string{"AES_ENCODE", "AES_DECODE"})
	if err != nil {
		panic(err)
	}
	return mode
}

func needBase64DecodeEncode() bool {
	needBase64Decode, err := dlgs.Question("二次解密/加密", "是否需要base64二次解密/加密", false)
	if err != nil {
		panic(err)
	}
	return needBase64Decode
}

func needContinue() bool {
	needCont, err := dlgs.Question("是否继续操作", "是否继续操作", false)
	if err != nil {
		panic(err)
	}
	return needCont
}
func checkClose() {
	needClose, err := dlgs.Question("是否关闭", "是否关闭AES工具", false)
	if err != nil {
		panic(err)
	}
	if needClose {
		os.Exit(0)
	}
}

func validLogin(name, pass string) bool {
	err := new(User).Login(name, pass)
	if err != nil {
		_, err = new(Log).Insert(Log{Operator: name, Content: "FAIL", Operation: "LOGIN"})
		return false
	}
	if err != nil {
		_, err = new(Log).Insert(Log{Operator: name, Content: "FAIL: " + err.Error(), Operation: "LOGIN"})
		return false
	}
	return true
}

func recoverTop() {
	if r := recover(); r != nil {
		PrintError(fmt.Sprintf("AES ERROR %v", r))
	}
}
