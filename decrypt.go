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
	"sync"
	"sync/atomic"
	"time"
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
		needBase64Bool = needBase64DecodeEncode()
		if strings.Contains(filePath, ".xls") {
			decodeExcel(filePath)
		} else if strings.Contains(filePath, ".txt") {
			decodeTxt(filePath)
		} else {
			PrintError("当前仅支持 txt/xls/xlsx 格式的文件")
		}
	}
}

func decodeTxt(filePath string) {
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
		cnt  int
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

var (
	eachTimes        = 40000 //每40000个处理一次save
	goNum            = 3     //开启每次处理协程数量
	excelFile        *excelize.File
	needBase64Bool   bool
	successCnt       int32
	dataColumnNum    int
	insertColumnName string
)

func decodeExcel(filePath string) {
	var (
		errW, dataColumnName string
		err                  error
	)
	dataColumnName, _, err = dlgs.List("模式", "请选择解密数据所在列:",
		[]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L"})
	dataColumnNum = ColumnNameToNumber(dataColumnName)
	insertColumnName = ColumnNumberToName(dataColumnNum + 1) //插入新列用于存放手机号

	excelFile, err = excelize.OpenFile(filePath)
	if err != nil {
		PrintError("加载文件时出错 " + err.Error() + "\n你可以尝试新建或另存为.xlsx文件")
	}
	defer excelFile.Save()
	sheets := excelFile.GetSheetMap()
	for _, sheet := range sheets {
		excelFile.InsertCol(sheet, insertColumnName) //需要解密数据列后插入新列
		rows := excelFile.GetRows(sheet)
		PrintMessage("reader file success,now start processing...")
		if err == nil {
			total := len(rows)
			if total > eachTimes { //文件超过3w并行处理
				segmens := SplitArrayStep(rows, eachTimes)
				for _, segmen := range segmens {
					var ws = sync.WaitGroup{}
					if total < segmen+eachTimes {
						process(sheet, rows[segmen:], segmen, nil)
					} else {
						for i := 0; i < goNum; i++ {
							ws.Add(1)
							start := segmen + i*eachTimes/goNum
							go process(sheet, rows[start:start+eachTimes/goNum], start, &ws)
						}
						ws.Wait()
					}
				}
			} else {
				process(sheet, rows, 0, nil)
			}
		}
	}
	if errW != "" {
		logger(fmt.Sprintf("FAIL: file: %s,err: %s", filePath, errW), "DECODE_EXCEL_ERROR")
	}
	logger(fmt.Sprintf("SUCCESS: total: %d", successCnt), "DECODE_EXCEL_SUCCESS")
	Notify(fmt.Sprintf(" SUCCESS total: %d", successCnt))

}

func process(sheet string, rows [][]string, start int, ws *sync.WaitGroup) {
	defer func(ws *sync.WaitGroup) {
		if ws != nil {
			ws.Done()
		}
	}(ws)
	var errW string
	var localSuccess int
	dataColumnName := ColumnNumberToName(dataColumnNum)
	for index, row := range rows {
		content, err := DecodeAes(row[dataColumnNum])
		if err != nil {
			errW += fmt.Sprintf("部分信息解密Aes失败，位置为 %s%d, 原始数据为 %v，错误为 %v \n",
				dataColumnName, index+1, row[dataColumnNum], err.Error())
			continue
		}
		if needBase64Bool {
			content, err = DecodeBase64(content)
			if err != nil {
				errW += fmt.Sprintf("部分信息Base64失败，位置为 %s%d,原始数据为 %v，错误为 %v \n",
					dataColumnName, index+1, row[dataColumnNum], err.Error())
				continue
			}
		}
		localSuccess++
		excelFile.SetCellStr(sheet, insertColumnName+strconv.Itoa(index+1+start), content)
	}
	atomic.AddInt32(&successCnt, int32(localSuccess))
	if errW != "" {
		logger("FAIL: when process: "+errW, "DECODE_EXCEL_PROCESS_FAIL")
	}
	Notify(fmt.Sprintf("Aes: already finished cnt: %d", successCnt))
	logger(fmt.Sprintf("success process: %v,%v,%v", start, localSuccess, time.Now()), "DECODE_EXCEL_PROCESS")
}
