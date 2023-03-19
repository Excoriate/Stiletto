package display

import (
	"fmt"
	"github.com/pterm/pterm"
	"strings"
)

func UXError(title string, msg string, err error) {
	pterm.Error.Prefix = pterm.Prefix{
		Text:  strings.ToUpper(title),
		Style: pterm.NewStyle(pterm.BgCyan, pterm.FgRed),
	}

	var errMsg string
	if err != nil {
		if msg != "" {
			errMsg = fmt.Sprintf("%s: %s", msg, err)
		} else {
			errMsg = err.Error()
		}
	}

	pterm.Error.Println(errMsg)
}

func UXInfo(title string, msg string) {
	pterm.Info.Prefix = pterm.Prefix{
		Text:  strings.ToUpper(title),
		Style: pterm.NewStyle(pterm.BgCyan, pterm.FgBlack),
	}
	pterm.Info.Println(msg)
}

func UXSuccess(title string, msg string) {
	pterm.Success.Prefix = pterm.Prefix{
		Text:  strings.ToUpper(title),
		Style: pterm.NewStyle(pterm.BgCyan, pterm.FgBlack),
	}
	pterm.Success.Println(msg)
}

func UXWarning(title string, msg string) {
	pterm.Warning.Prefix = pterm.Prefix{
		Text:  strings.ToUpper(title),
		Style: pterm.NewStyle(pterm.BgCyan, pterm.FgBlack),
	}
	pterm.Warning.Println(msg)
}
