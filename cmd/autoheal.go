package main

import (
	"github.com/navidnabavi/auto-healer/internal/autoheal"
)

func main() {
	autoHealer := autoheal.NewAutoHealer()
	autoHealer.Spin()
}
