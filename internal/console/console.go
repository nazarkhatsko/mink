package console

import "fmt"

func Print(format string, args ...any) {
	fmt.Printf(format, args...)
}

func Println(args ...any) {
	fmt.Println(args...)
}

func Printlnf(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}
