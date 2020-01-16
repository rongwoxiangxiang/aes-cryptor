package main

func encodeAESString(strArr []string) (data map[string]string) {
	if len(strArr) > 0 {
		data = make(map[string]string, len(strArr))
		for _, str := range strArr {
			data[str] = EncodeAes(str)
		}
	}
	return data
}

func encodeBase64String(strArr []string) (data map[string]string) {
	if len(strArr) > 0 {
		data = make(map[string]string, len(strArr))
		for _, str := range strArr {
			data[str] = EncodeBase64(str)
		}
	}
	return data
}

func encodeBase64AesString(strArr []string) (data map[string]string) {
	if len(strArr) > 0 {
		data = make(map[string]string, len(strArr))
		for _, str := range strArr {
			base64encodeStr := EncodeBase64(str)
			data[str] = EncodeAes(base64encodeStr)
		}
	}
	return data
}
