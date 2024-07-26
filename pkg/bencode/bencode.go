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

// Decoded result linked list node
type DecodedCodeNode[T any] struct {
	Val        T
	Next, Tail *DecodedCodeNode[T]
}

// Push value to the end of the decoded list
func (node *DecodedCodeNode[T]) Push(val T) {
	tail := node.Tail

	if tail == nil {
		node.Val = val
		tail = node
	} else {
		tail.Next = &DecodedCodeNode[T]{
			Val: val,
		}
		node.Tail = tail.Next
	}
}

// Decode bencode to integer
func decodeInt(bCode *string) (int, int, error) {
	bLen := len(*bCode)

	// Finding "e" index
	var index int
	for i := 0; i < bLen; i++ {
		curr := string((*bCode)[i])
		if curr == "e" {
			index = i
			break
		}
	}

	// Remove "i" and "e" from string
	intStr := (*bCode)[1:index]

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

	intLen, err := strconv.Atoi(strLen)
	if err != nil {
		panic("Could not parse string length")
	}

	index := start + intLen
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
	// res, err := strconv.Atoi(intStr)
	// if err != nil {
	// 	panic("Invalid bencode")
	// }

	return string(intStr), index, nil
}

// Decode bencode to dictionary -> tree
func decodeDict(bCode *string) (string, int, error) {

	return "", -1, nil
}

// Decode bencode string
func Decode(bCode *string) (*DecodedCodeNode[any], error) {
	bLen := len(*bCode)
	res := &DecodedCodeNode[any]{}

	for i := 0; i < bLen; i++ {
		fChar := string((*bCode)[i])
		remain := (*bCode)[i:bLen]

		fmt.Println(fChar)
		fmt.Println(remain)

		switch fChar {
		case "i":
			temp, end, err := decodeInt(&remain)
			if err != nil {
				panic(fmt.Sprintf("[Decoding Error] %v\n", err))
			}

			i = end
			res.Push(temp)
			continue
		case "l":
			temp, end, err := decodeInt(&remain)
			if err != nil {
				panic(fmt.Sprintf("[Decoding Error] %v\n", err))
			}

			i = end
			res.Push(temp)
			continue
		case "d":
			temp, end, err := decodeDict(&remain)
			if err != nil {
				panic(fmt.Sprintf("[Decoding Error] %v\n", err))
			}

			i = end
			res.Push(temp)
			continue
		default:
			if _, err := strconv.Atoi(fChar); err == nil {
				for j := i; j < bLen; j++ {
					curr := string((*bCode)[j])

					if curr == ":" {
						temp, end, err := decodeString(&remain)
						if err != nil {
							panic("[Decode error] Invalid bencode")
						}

						i = end
						res.Push(temp)
						break
					} else if _, err := strconv.Atoi(curr); err == nil {
						continue
					} else {
						panic("[Decode error] Invalid bencode")
					}
				}
			} else {
				panic("[Decode error] Invalid bencode")
			}
			break
		}
	}

	return res, nil
}
