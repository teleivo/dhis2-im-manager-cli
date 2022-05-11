package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func main() {
	blockA := "AAA\nAAA\nAAA\nAAA\nAAA"
	blockB := "BBB\nBBB\nBBB"
	str := lipgloss.JoinHorizontal(0.2, blockA, blockB)
	fmt.Println(str)
}
