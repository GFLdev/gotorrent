/*
 * bencode.go
 * Author: Gabriel Franco Lourenco
 * Github: github.com/GFLdev
 */

package bencode

import (
	"fmt"
	"strconv"
)

// Decode bencode to integer
func decodeInt(bCode *string) (int, int, error) {
	bLen := len(*bCode)

	// Finding "e" index
	var index int
	for i := 0; i < bLen-1; i++ {
		index = i
		if (*bCode)[i] == 'e' {
			break
		}
	}

	// Remove "i" and "e" from string
	intStr := (*bCode)[1 : index-1]

	// Try to convert string to integer
	res, err := strconv.Atoi(intStr)
	if err != nil {
		panic("Invalid bencode")
	}

	return res, index, nil
}

// Decode bencode to string
func decodeString(bCode *string) (string, int, error) {
	// Find string length
	strLen := ""
	start := 0
	for i := 0; i < len(*bCode); i++ {
		if (*bCode)[i] == ':' {
			start = i + 1
			break
		}
		strLen += string((*bCode)[i])
	}

	len, err := strconv.Atoi(strLen)
	if err != nil {
		panic("Could not parse string length")
	}

	index := start + len
	return (*bCode)[start:index], index, nil
}

// Decode bencode to list -> array
func decodeList(bCode *string) (string, int, error) {
	bLen := len(*bCode)

	// Finding "e" index
	var index int
	extra := 0
	for i := 0; i < bLen-1; i++ {
		if (*bCode)[i] == 'l' || (*bCode)[i] == 'd' || (*bCode)[i] == 'i' {
			extra++
		} else if (*bCode)[i] == 'e' {
			extra--
			if extra == -1 {
				index = i
				break
			}
		}
	}

	// Remove "l" and "e" from string
	intStr := (*bCode)[1 : index-1]

	// Try to convert string to integer
	res, err := strconv.Atoi(intStr)
	if err != nil {
		panic("Invalid bencode")
	}

	return string(res), index, nil
}

// Decode bencode to dictionary -> tree
func decodeDict(bCode *string) (string, int, error) {

	return "", -1, nil
}

// Decode bencode string
func Decode(bCode string) (string, error) {
	fChar := bCode[0]
	bLen := len(bCode)
	res := ""

	for i := 0; i < bLen; i++ {
		if fChar == 'i' {
			temp, end, err := decodeInt(&bCode)
			if err != nil {
				panic(fmt.Sprintf("[Decoding Error] %v\n", err))
			}

			i = end
			res += string(temp) + "\n"
		} else if fChar == 'l' {
			res, end, err := decodeList(&bCode)
			if err != nil {
				panic(fmt.Sprintf("[Decoding Error] %v\n", err))
			}

			i = end
			res += string(res) + "\n"
		} else if fChar == 'd' {
			temp, end, err := decodeDict(&bCode)
			if err != nil {
				panic(fmt.Sprintf("[Decoding Error] %v\n", err))
			}

			i = end
			res += string(temp) + "\n"
		} else if _, err := strconv.Atoi(string(fChar)); err == nil {
			for j := 0; j < bLen; j++ {
				_, err := strconv.Atoi(string(bCode[i]))

				if bCode[j] == ':' && err == nil {
					temp, end, err := decodeString(&bCode)
					if err != nil {
						panic("[Decode error] Invalid bencode")
					}

					i = end
					res = string(temp) + "\n"
					continue
				} else {
					panic("[Decode error] Invalid bencode")
				}
			}
		} else {
			panic("[Decode error] Invalid bencode")
		}
	}

	panic("[Decode error] Invalid bencode")
}
