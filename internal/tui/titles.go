package tui

import (
	"fmt"
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"strings"
)

type TUITitle struct {
}

func (t *TUITitle) ShowTitleAndDescription(title, description string) {
	titleNormalised := strings.TrimSpace(strings.ToUpper(title))
	s, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString(titleNormalised)).
		Srender()
	pterm.DefaultCenter.Println(s)

	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(description)
}

func (t *TUITitle) ShowTitle(title string) {
	titleNormalised := strings.TrimSpace(strings.ToUpper(title))
	s, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString(titleNormalised)).
		Srender()
	pterm.DefaultCenter.Println(s)
}

func (t *TUITitle) ShowDescription(description string) {
	subtitleNormalised := strings.TrimSpace(strings.ToUpper(description))
	pterm.Println()
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println("--------------------------------")
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(subtitleNormalised)
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println("--------------------------------")
	pterm.Println()
}

func (t *TUITitle) ShowInitDetails(jobName, taskName, workDir, mountDir, targetDir string) {
	pterm.Println()
	pterm.DefaultBasicText.Println(pterm.LightWhite("--------------------------------------------------"))
	pterm.DefaultBasicText.Println("Job" + pterm.LightMagenta(fmt.Sprintf(" %s ",
		common.NormaliseStringUpper(jobName))))

	pterm.DefaultBasicText.Println("TaskName ->" + pterm.LightMagenta(fmt.Sprintf(" %s ",
		common.NormaliseStringUpper(taskName))))

	pterm.DefaultBasicText.Println("Workdir ->" + pterm.LightMagenta(fmt.Sprintf(" %s ",
		common.NormaliseStringUpper(workDir))))

	pterm.DefaultBasicText.Println("MountDir ->" + pterm.LightMagenta(fmt.Sprintf(" %s ",
		common.NormaliseStringUpper(mountDir))))

	pterm.DefaultBasicText.Println("TargetDir ->" + pterm.LightMagenta(fmt.Sprintf(" %s ",
		common.NormaliseStringUpper(targetDir))))

	pterm.DefaultBasicText.Println(pterm.LightWhite("--------------------------------------------------"))
	pterm.Println()
	pterm.Println()
}

func (t *TUITitle) ShowTaskDetails(taskName, actionName, workDir, mountDir, targetDir string) {
	pterm.Println()
	pterm.DefaultBasicText.Println(pterm.LightWhite("--------------------------------------------------"))
	pterm.DefaultBasicText.Println("TaskName" + pterm.LightMagenta(fmt.Sprintf(" %s ",
		common.NormaliseStringUpper(taskName))))

	pterm.DefaultBasicText.Println("actionName ->" + pterm.LightMagenta(fmt.Sprintf(" %s ",
		common.NormaliseStringUpper(actionName))))

	pterm.DefaultBasicText.Println("Workdir ->" + pterm.LightMagenta(fmt.Sprintf(" %s ",
		common.NormaliseStringUpper(workDir))))

	pterm.DefaultBasicText.Println("MountDir ->" + pterm.LightMagenta(fmt.Sprintf(" %s ",
		common.NormaliseStringUpper(mountDir))))

	pterm.DefaultBasicText.Println("TargetDir ->" + pterm.LightMagenta(fmt.Sprintf(" %s ",
		common.NormaliseStringUpper(targetDir))))

	pterm.DefaultBasicText.Println(pterm.LightWhite("--------------------------------------------------"))
	pterm.Println()
	pterm.Println()
}

func (t *TUITitle) ShowSubTitle(mainTitle string, subTitle string) {
	_ = pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle(strings.ToUpper(mainTitle), pterm.NewStyle(pterm.FgCyan)),
		putils.LettersFromStringWithStyle(strings.ToUpper(subTitle), pterm.NewStyle(pterm.FgLightMagenta))).
		Render()
}

func NewTitle() TUIDisplayer {
	return &TUITitle{}
}
