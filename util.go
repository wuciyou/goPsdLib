package goPsdLib

func UnpackRLEBits(data []byte, length int) []byte {
	result := make([]byte, length)
	wPos, rPos := 0, 0
	for rPos < len(data) {
		n := data[rPos]
		rPos++
		if n < 128 {
			count := int(n) + 1
			for j := 0; j < count; j++ {
				// log.Printf("result[%d] = data[%d], datalen(%d) \n ", wPos, rPos, len(data))
				result[wPos] = byte(data[rPos])
				wPos++
				rPos++
			}
		} else {
			b := byte(data[rPos])
			rPos++
			count := int(-n) + 1
			for j := 0; j < count; j++ {
				result[wPos] = b
				wPos++
			}
		}
	}
	return result
}

func StringValueIs(value string, values ...string) bool {
	for i := range values {
		if value == values[i] {
			return true
		}
	}
	return false
}
