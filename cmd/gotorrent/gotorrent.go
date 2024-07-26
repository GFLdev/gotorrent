/*
 * bencode.go
 * Author: Gabriel Franco Lourenco
 * Github: github.com/GFLdev
 */

package main

import (
	"fmt"
	"os"

	"github.com/GFLdev/gotorrent/pkg/bencode"
)

func main() {
	bCode := os.Args[1]
	res, _ := bencode.Decode(&bCode)
	for res != nil {
		fmt.Println(res.Val)
		res = res.Next
	}
}
