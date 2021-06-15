package example

import (
	"fmt"
	awwrs "go-awrrs/awrrs"
	"testing"
)

func TestWrrSlice_Add(t *testing.T) {
	wrrs := &awwrs.WrrSlice{}
	wrrs.Add("hello", 5)
	wrrs.Add("say", 4)
	wrrs.Add("cool", 6)
	for i := 0; i < 10; i++ {
		next, err := wrrs.Next()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("call func：", next)
	}
}
