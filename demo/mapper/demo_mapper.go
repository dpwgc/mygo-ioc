package mapper

import "fmt"

type DemoMapper struct {
}

func (m *DemoMapper) Print(text string) {
	fmt.Println("print:", text)
}
