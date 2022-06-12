package color

func ANSI_COLOR(color_name string) string {
	//colorReset := "\033[0m"
	switch color_name {
	case "RED":
		colorRed := "\033[31m"
		return colorRed
	case "GREEN":
		colorGreen := "\033[32m"
		return colorGreen
	case "RESET":
		colorReset := "\033[0m"
		return colorReset
	}

	/* colorGreen := "\033[32m"
	colorYellow := "\033[33m"
	colorBlue := "\033[34m"
	colorPurple := "\033[35m"
	colorCyan := "\033[36m"
	colorWhite := "\033[37m" */

	/* fmt.Println(string(colorRed), "test")
	fmt.Println(string(colorGreen), "test")
	fmt.Println(string(colorYellow), "test")
	fmt.Println(string(colorBlue), "test")
	fmt.Println(string(colorPurple), "test")
	fmt.Println(string(colorWhite), "test")
	fmt.Println(string(colorCyan), "test", string(colorReset))
	fmt.Println("next") */
	colorWhite := "\033[37m"
	return colorWhite
}
