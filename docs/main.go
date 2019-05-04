package main

import (
	"fmt"
	"github.com/huntsman-li/resp"
)

func main() {
	//decoder()
	encoder()
}

func decoder() {
	var dest string
	if err := resp.Unmarshal([]byte("$3\r\nFoo\r\n"), &dest); err != nil {
		panic(err)
	}
	fmt.Println(dest)
}

func encoder() {
	buf, err := resp.Marshal("Foo") // RESP: $3\r\nFoo\r\n
	if err != nil {
		panic(err)
	}

	fmt.Printf("buf: %s\n", string(buf))
}
