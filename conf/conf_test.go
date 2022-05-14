package conf

import (
	"fmt"
	"testing"
)

func Test_GetCoinList(t *testing.T) {
	err := InitViper()
	if err != nil {
		panic(fmt.Sprintf("Fatal error config file: <%v> ", err))
	}
	coinList := GetCoinList()
	fmt.Println(coinList)
}
