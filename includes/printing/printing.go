package printing

import "fmt"

const (
	C_GRAY    = "\x1b[90m"
	C_RED     = "\x1b[91m"
	C_GREEN   = "\x1b[92m"
	C_YELLOW  = "\x1b[93m"
	C_BLUE    = "\x1b[94m"
	C_MAGENTA = "\x1b[95m"
	C_CYAN    = "\x1b[96m"
	C_CYAN2   = "\x1b[36m"

	C_BG_RED   = "\x1b[41m"
	C_BG_GREEN = "\x1b[42m"
	C_BG_BLUE  = "\x1b[44m"
	C_BG_YELL  = "\x1b[43m"
	C_BG_MAG   = "\x1b[45m"
	C_BG_CYAN  = "\x1b[46m"

	C_BOLD    = "\x1b[1m"
	C_PALE    = "\x1b[2m"
	C_INTALIC = "\x1b[3m"
	C_UNDER   = "\x1b[4m"
	C_FLASH   = "\x1b[5m"

	C_RST = "\x1b[0m"
)

var PRINT_ON = true // add ability to turn off printf globally

func Pf(format string, va ...any) {
	//return
	if PRINT_ON {
		fmt.Printf(format, va...)
	}
}

func PfGray(fmt string, va ...any)    { Pf(C_GRAY+fmt+C_RST, va...) }
func PfRed(fmt string, va ...any)     { Pf(C_RED+fmt+C_RST, va...) }
func PfGreen(fmt string, va ...any)   { Pf(C_GREEN+fmt+C_RST, va...) }
func PfBlue(fmt string, va ...any)    { Pf(C_BLUE+fmt+C_RST, va...) }
func PfMagenta(fmt string, va ...any) { Pf(C_MAGENTA+fmt+C_RST, va...) }
func PfYellow(fmt string, va ...any)  { Pf(C_YELLOW+fmt+C_RST, va...) }
func PfBold(fmt string, va ...any)    { Pf(C_BOLD+fmt+C_RST, va...) }
