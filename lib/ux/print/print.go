package print

import (
	"fmt"

	. "github.com/christiangelone/bang/lib/sugar"
	"github.com/fatih/color"
)

const (
	FgNoColor = -1
	FgWhite   = color.FgWhite
	FgYellow  = color.FgHiYellow
	FgRed     = color.FgHiRed
	FgGreen   = color.FgHiGreen
	FgBlue    = color.FgHiBlue
	Padding   = "    "
)

type FgColor = color.Attribute

func Sprint(seq ...interface{}) string {
	var currentColor FgColor = FgNoColor
	var text = ""
	for _, e := range seq {
		switch this := e.(type) {
		case FgColor:
			currentColor = this
		case string:
			if text == "" {
				text = If(currentColor == FgNoColor)(this).Else(color.New(currentColor).Sprint(this)).(string)
			} else {
				text = fmt.Sprintf("%s %s", text, If(currentColor == FgNoColor)(this).Else(color.New(currentColor).Sprint(this)).(string))
			}
		default:
		}
	}

	return text
}

func Print(seq ...interface{}) {
	var text = Sprint(seq...)
	fmt.Print(text)
}

func Println(seq ...interface{}) {
	var text = Sprint(seq...)
	fmt.Println(text)
}

func Bullet(str string) {
	Print(FgBlue, Padding+"• ")
	Println(FgNoColor, str)
}
