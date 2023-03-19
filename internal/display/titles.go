package display

import (
	"github.com/pterm/pterm"
	"strings"
)

func UXTitleAndDescription(title string, description string) {
	titleNormalised := strings.TrimSpace(strings.ToUpper(title))
	s, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString(titleNormalised)).
		Srender()
	pterm.DefaultCenter.Println(s)

	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(description)
}

func UXTitle(title string) {
	titleNormalised := strings.TrimSpace(strings.ToUpper(title))
	s, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString(titleNormalised)).
		Srender()
	pterm.DefaultCenter.Println(s)
}

func UXSubTitle(subtitle string) {
	subtitleNormalised := strings.TrimSpace(strings.ToUpper(subtitle))
	pterm.Println()
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println("--------------------------------")
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(subtitleNormalised)
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println("--------------------------------")
	pterm.Println()
}
