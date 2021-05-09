package packet

import (
	"bytes"
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	var buf = bytes.NewBuffer(nil)
	buf.WriteString("123321321")
	fmt.Println("buff=", buf.Len())

	buf.Next(4)
	fmt.Println("buff2=", buf.Len())
	fmt.Println(string(buf.Bytes()))
	if buf.Len() == 0 {
		buf.Reset()
	}
}
