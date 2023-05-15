package config

import (
	"github.com/Excoriate/stiletto/internal/errors"
	"github.com/Excoriate/stiletto/internal/tui"
	"github.com/spf13/viper"
)

type CLIGlobalArgs struct {
	WorkingDir                     string
	MountDir                       string
	TargetDir                      string
	TaskName                       string
	ScanEnvVarKeys                 []string
	EnvKeyValuePairsToSet          map[string]interface{}
	EnvKeyValuePairsToSetString    map[string]string
	ScanAWSKeys                    bool
	ScanTerraformVars              bool
	ScanAllEnvVars                 bool
	CustomCommands                 []string
	InitDaggerWithWorkDirByDefault bool
	RunInVendor                    bool
}

func GetCLIGlobalArgs() (CLIGlobalArgs, error) {
	cfg := Cfg{}
	defaultEmptyMap := make(map[string]interface{})

	// 'set-env' option
	keValuePairsFromViper, err := cfg.GetFromViperOrDefault("set-env", defaultEmptyMap)
	if err != nil {
		return CLIGlobalArgs{}, errors.NewArgumentError(
			"Error trying to parse or resolve argument 'set-env'", err)
	}

	setEnvValue := keValuePairsFromViper.Value.(map[string]interface{})

	// 'scan-env' option
	defaultEmptySliceString := make([]string, 0)
	scanEnvVarKeysFromViper, err := cfg.GetFromViperOrDefault("scan-env", defaultEmptySliceString)
	if err != nil {
		return CLIGlobalArgs{}, errors.NewArgumentError(
			"Error trying to parse or resolve argument 'scan-env'", err)
	}

	scanEnvVarKeys := scanEnvVarKeysFromViper.Value.([]string)

	args := CLIGlobalArgs{
		WorkingDir:            viper.Get("work-dir").(string),
		MountDir:              viper.Get("mount-dir").(string),
		TargetDir:             viper.Get("target-dir").(string),
		TaskName:              viper.Get("task").(string),
		ScanEnvVarKeys:        scanEnvVarKeys,
		EnvKeyValuePairsToSet: setEnvValue,
		ScanAWSKeys:           viper.Get("scan-aws-keys").(bool),
		ScanTerraformVars:     viper.Get("scan-terraform-vars").(bool),
		ScanAllEnvVars:        viper.Get("scan-all-env-vars").(bool),
		//CustomCommands:                 viper.Get("custom-cmds").([]string),
		CustomCommands:                 []string{},
		InitDaggerWithWorkDirByDefault: viper.Get("init-dagger-with-workdir").(bool),
		RunInVendor:                    viper.Get("run-in-vendor").(bool),
	}

	for k, v := range args.EnvKeyValuePairsToSet {
		args.EnvKeyValuePairsToSetString[k] = v.(string)
	}

	return args, nil
}

func ShowCLITitle() {
	ux := tui.TUITitle{}
	ux.ShowTitleAndDescription("STILETTO",
		"Stiletto is a pipeline framework that works on top of Dagger.io. "+
			"Makes your pipelines more readable and easier to maintain.")
}
