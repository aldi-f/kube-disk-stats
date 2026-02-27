package display

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	colorGreen  = color.New(color.FgGreen).SprintFunc()
	colorYellow = color.New(color.FgYellow).SprintFunc()
	colorRed    = color.New(color.FgRed).SprintFunc()
)

func ColorizePercentage(percentage float64) string {
	if percentage < 60 {
		return colorGreen(percentage)
	} else if percentage < 80 {
		return colorYellow(percentage)
	}
	return colorRed(percentage)
}

func ColorizePercentageText(percentage float64) string {
	text := fmt.Sprintf("%.1f%%", percentage)
	if percentage < 60 {
		return colorGreen(text)
	} else if percentage < 80 {
		return colorYellow(text)
	}
	return colorRed(text)
}

func GetColorForPercentage(percentage float64) func(a ...interface{}) string {
	if percentage < 60 {
		return colorGreen
	} else if percentage < 80 {
		return colorYellow
	}
	return colorRed
}
