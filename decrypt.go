package main

import (
	"bufio"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gen2brain/dlgs"
	"io"
	"os"
	"strconv"
	"strings"
)

func decodeAesString(strArr []string) (datas map[string]string, errW string) {
	if len(strArr) > 0 {
		datas = make(map[string]string, len(strArr))
		for _, str := range strArr {
			data, err := DecodeAes(str)
			if err != nil {
				errW += fmt.Sprintf("解密Aes时出错，原始数据为 %v， 错误为 %v \n", str, err.Error())
				continue
			}
			datas[str] = data
		}
	}
	return
}

func decodeAesBase64String(strArr []string) (datas map[string]string, errW string) {
	if len(strArr) > 0 {
		datas = make(map[string]string, len(strArr))
		for _, str := range strArr {
			data, err := DecodeAes(str)
			if err != nil {
				errW += fmt.Sprintf("解密Aes时出错，原始数据为 %v， 错误为 %v \n", str, err.Error())
				continue
			}
			data, err = DecodeBase64(data)
			if err != nil {
				errW += fmt.Sprintf("解密Base64时出错，原始数据为 %v， 错误为 %v \n", str, err.Error())
				continue
			}
			datas[str] = data
		}
	}
	return
}

func decodeBase64String(strArr []string) (datas map[string]string, errW string) {
	if len(strArr) > 0 {
		datas = make(map[string]string, len(strArr))
		for _, str := range strArr {
			data, err := DecodeBase64(str)
			if err != nil {
				errW += fmt.Sprintf("解密Base64时出错，原始数据为 %v， 错误为 %v \n", str, err.Error())
				continue
			}
			datas[str] = data
		}
	}
	if errW != "" {
		PrintWarning(errW)
	}
	return
}

func decodeFile() {
	filePath := getFile()
	_, err := os.Stat(filePath)
	if err == nil {
		if strings.Contains(filePath, ".xls") {
			decodeExcel(filePath, needBase64DecodeEncode())
		} else if strings.Contains(filePath, ".txt") {
			decodeTxt(filePath, needBase64DecodeEncode())
		} else {
			PrintError("当前仅支持 txt/xls/xlsx 格式的文件")
		}
	}
}

func decodeTxt(filePath string, needBase64Bool bool) {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		PrintError("【0】打开源文件出错：" + err.Error())
	}
	file2, err := os.OpenFile(filePath+"-decode.txt", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		PrintError("【1】创建解密文件出错：" + err.Error())
	}
	var (
		errW string
		cnt int
	)
	PrintMessage("处理中，请稍后...")
	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if errW != "" {
					PrintWarning(errW)
					logger(fmt.Sprintf("FAIL: file: %s,err: %s", filePath, errW), "DECODE_TEXT_ERROR")
				}
				logger(fmt.Sprintf("SUCCESS: file: %s,cnt: %d", filePath, cnt), "DECODE_TEXT_SUCCESS")
				PrintSuccess("SUCCESS! \n\n解密文件为" + file2.Name())
			} else {
				PrintWarning("【end】读入加密文件出错! " + err.Error())
			}
			os.Exit(0)
		}
		line = strings.TrimSpace(line)
		content, err := DecodeAes(line)
		if err != nil {
			errW += fmt.Sprintf("解密Aes时出错，原始数据为 %v，错误为 %v \n", line, err.Error())
			continue
		}
		if needBase64Bool {
			content, err = DecodeBase64(content)
			if err != nil {
				errW += fmt.Sprintf("解密Base64时出错，原始数据为 %v，错误为 %v \n", line, err.Error())
				continue
			}
		}
		_, err = file2.WriteString(content + "\n")
		if err != nil {
			errW += fmt.Sprintf("写入时出错，原始数据为 %v, 解密后为 %v，错误为 %v \n", line, content, err.Error())
			continue
		}
		cnt++
	}
}

func decodeExcel(filePath string, needBase64Bool bool) {
	var (
		errW string
	 	errNums int
		cnt int
	)
	columnName, _, err := dlgs.List("模式", "请选择解密数据所在列:",
		[]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L"})
	columnNum := ColumnNameToNumber(columnName)
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		PrintError("加载文件时出错 " + err.Error() + "\n你可以尝试新建或另存为.xlsx文件")
	}
	PrintMessage("处理中，请稍后...")
	sheets := f.GetSheetMap()
	for _, sheet := range sheets {
		rows := f.GetRows(sheet)
		if err == nil {
			col := ColumnNumberToName(len(rows[0]) + 2)
			for index, row := range rows {
				if errNums > 30 {
					errW = "错误数据过多，请修改后重新解密\n" + errW
					break
				}
				content, err := DecodeAes(row[columnNum])
				if err != nil {
					errNums++
					errW += fmt.Sprintf("部分信息解密Aes失败，位置为 %s%d, 原始数据为 %v，错误为 %v \n",
						columnName, index+1, row[columnNum], err.Error())
					continue
				}
				if needBase64Bool {
					content, err = DecodeBase64(content)
					if err != nil {
						errNums++
						errW += fmt.Sprintf("部分信息Base64失败，位置为 %s%d,原始数据为 %v，错误为 %v \n",
							columnName, index+1, row[columnNum], err.Error())
						continue
					}
				}
				cnt++
				f.SetCellStr(sheet, col+strconv.Itoa(index+1), content)
			}
		}
	}
	if err := f.Save(); err != nil {
		PrintError("保存文件时出错 " + err.Error())
	}
	if errW != "" {
		logger(fmt.Sprintf("FAIL: file: %s,err: %s", filePath, errW), "DECODE_EXCEL_ERROR")
	}
	logger(fmt.Sprintf("SUCCESS: file: %s,cnt: %d", filePath, cnt), "DECODE_EXCEL_SUCCESS")
	PrintSuccess(fmt.Sprintf("SUCCESS \nSUCCESS \nSUCCESS \n\n成功数量: %d, err: %v", cnt, errW))
	ch<-true
	return
}
