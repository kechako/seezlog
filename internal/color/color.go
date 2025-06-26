package color

type Color string

const (
	Reset     Color = "\x1b[0m"
	Bold      Color = "\x1b[1m"
	Dim       Color = "\x1b[2m"
	Italic    Color = "\x1b[3m"
	Underline Color = "\x1b[4m"

	Black   Color = "\x1b[30m"
	Red     Color = "\x1b[31m"
	Green   Color = "\x1b[32m"
	Yellow  Color = "\x1b[33m"
	Blue    Color = "\x1b[34m"
	Magenta Color = "\x1b[35m"
	Cyan    Color = "\x1b[36m"
	White   Color = "\x1b[37m"

	BrightBlack   Color = "\x1b[90m"
	BrightRed     Color = "\x1b[91m"
	BrightGreen   Color = "\x1b[92m"
	BrightYellow  Color = "\x1b[93m"
	BrightBlue    Color = "\x1b[94m"
	BrightMagenta Color = "\x1b[95m"
	BrightCyan    Color = "\x1b[96m"
	BrightWhite   Color = "\x1b[97m"
)
