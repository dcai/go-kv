package main

const END = "\033[0m"

func FgGray(str string) string {
	return "\033[1;30m" + str + END
}

func FgRed(str string) string {
	return "\033[0;31m" + str + END
}

func FgYellow(str string) string {
	return "\033[1;33m" + str + END
}

func FgGreen(str string) string {
	return "\033[0;32m" + str + END
}
