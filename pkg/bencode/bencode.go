/*
 * bencode.go
 * Author: Gabriel Franco Lourenco
 * Github: github.com/GFLdev
 */

package bencode

import (
	"fmt"
)

// Decoded result linked list node
type DecodedCodeNode[T any] struct {
	Val        T
	Next, Tail *DecodedCodeNode[T]
}

// Get integer bencode string length
func getIntStrLen(bCode *string) int {
	bLen := len(*bCode)
	var index int

	// Finding "e" index
	for i := 0; i < bLen; i++ {
		curr := (*bCode)[i]
		if curr == 'e' {
			index = i
			break
		}
	}

	return index
}

// Get string bencode string starter point and length
func getStrStrLen(bCode *string) (int, int) {
	// Find string length
	start := 0
	for i := 0; i < len(*bCode); i++ {
		if (*bCode)[i] == ':' {
			start = i + 1
			break
		}
	}

	strLen := (*bCode)[0 : start-1]
	intLen := 0
	for _, byte := range strLen {
		if byte < '0' || byte > '9' {
			panic("Invalid bencode: could not parse string length")
		}

		intLen = intLen*10 + int(byte-'0')
	}

	index := start + intLen
	return start, index
}

// Get compound bencode string length, for lists and dictionaries
func getCompoundStrLen(bCode *string) int {
	bLen := len(*bCode)
	i := 0

	for ; i < bLen; i++ {
		switch (*bCode)[i] {
		case 'i':
			i += getIntStrLen(bCode)
			continue
		case 'l':
		case 'd':
			i += getCompoundStrLen(bCode)
			continue
		case 'e':
			break
		default:
			_, val := getStrStrLen(bCode)
			i += val
			continue
		}
	}

	return i
}

// Push value to the end of the decoded list
func (node *DecodedCodeNode[T]) Push(val T) {
	tail := node.Tail

	if tail == nil {
		node.Val = val
		node.Tail = node
	} else {
		tail.Next = &DecodedCodeNode[T]{
			Val: val,
		}
		node.Tail = tail.Next
	}
}

// Decode bencode to integer
func decodeInt(bCode *string) (int, int, error) {
	index := getIntStrLen(bCode)

	// Remove "i" and "e" from string
	intSlice := (*bCode)[1:index]

	// Convert bytes to integer
	res := 0
	for _, byte := range intSlice {
		if byte < '0' || byte > '9' {
			return 0, 0, fmt.Errorf("Invalid bencode: non-digit character found")
		}

		res = res*10 + int(byte-'0')
	}

	return res, index, nil
}

// Decode bencode to string
func decodeString(bCode *string) (string, int, error) {
	start, index := getStrStrLen(bCode)

	return (*bCode)[start:index], index - 1, nil
}

// Decode bencode to list
func decodeList(bCode *string) (*DecodedCodeNode[any], int, error) {
	bLen := len(*bCode)

	// Finding "e" index
	var index int
	extra := 0
	for i := 0; i < bLen; i++ {
		if (*bCode)[i] == 'l' || (*bCode)[i] == 'd' || (*bCode)[i] == 'i' {
			extra++
		} else if (*bCode)[i] == 'e' {
			extra--
			if extra == 0 {
				index = i
				break
			}
		}
	}

	// Remove "l" and "e" from string
	intStr := (*bCode)[1:index]

	// Convert list code to linked list
	res, err := Decode(&intStr)
	if err != nil {
		return nil, 0, fmt.Errorf("Invalid bencode: could not parse bencode list")
	}

	return res, index, nil
}

// Decode bencode to map
func decodeDict(bCode *string) (map[any]any, int, error) {
	bLen := len(*bCode)

	// Finding "e" index
	var index int
	extra := 0
	for i := 0; i < bLen; i++ {
		if (*bCode)[i] == 'l' || (*bCode)[i] == 'd' || (*bCode)[i] == 'i' {
			extra++
			fmt.Printf("Char: %v, Extra: %v, I: %v\n", string((*bCode)[i]), extra, i)
		} else if (*bCode)[i] == 'e' {
			extra--
			fmt.Printf("Char: %v, Extra: %v, I: %v\n", string((*bCode)[i]), extra, i)
			if extra == 0 {
				index = i
				break
			}
		}
	}

	// Remove "d" and "e" from string
	intStr := (*bCode)[1:index]
	fmt.Printf("intStr: %v\n", intStr)

	// Convert dictionary code to linked list
	temp, err := Decode(&intStr)
	if err != nil {
		return nil, 0, fmt.Errorf("Invalid bencode: could not parse bencode list")
	}

	// Convert linked list to JSON
	res := make(map[any]any)
	var key any
	for i := 1; ; i++ {
		if temp == nil {
			if i&1 == 0 {
				return nil, 0, fmt.Errorf("Invalid bencode: no value paired with key")
			}
			break
		}

		if i&1 == 1 {
			key = temp.Val
		} else {
			res[key] = temp.Val
		}

		temp = temp.Next
	}

	return res, index + 1, nil
}

// Decode bencode string
func Decode(bCode *string) (*DecodedCodeNode[any], error) {
	bLen := len(*bCode)
	res := &DecodedCodeNode[any]{}

	for i := 0; i < bLen; i++ {
		fChar := (*bCode)[i]
		remain := (*bCode)[i:bLen]
		fmt.Println(string(fChar))
		fmt.Println(remain)

		switch fChar {
		case 'i':
			temp, end, err := decodeInt(&remain)
			if err != nil {
				panic(fmt.Sprintf("[Decoding Error] %v\n", err))
			}

			i += end
			res.Push(temp)
			continue
		case 'l':
			temp, end, err := decodeList(&remain)
			if err != nil {
				panic(fmt.Sprintf("[Decoding Error] %v\n", err))
			}

			i += end
			res.Push(temp)
			continue
		case 'd':
			temp, end, err := decodeDict(&remain)
			if err != nil {
				panic(fmt.Sprintf("[Decoding Error] %v\n", err))
			}

			i += end
			res.Push(temp)
			continue
		default:
			for j := i; j < bLen; j++ {
				curr := (*bCode)[j]

				if curr == ':' {
					temp, end, err := decodeString(&remain)
					if err != nil {
						panic(fmt.Sprintf("[Decoding Error] %v\n", err))
					}

					i += end
					res.Push(temp)
					break
				} else if curr >= '0' && curr <= '9' {
					continue
				} else {
					panic("[Decoding Error] Invalid bencode")
				}
			}
			continue
		}
	}

	return res, nil
}
