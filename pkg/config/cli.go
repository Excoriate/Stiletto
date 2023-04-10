package config

import (
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
	CustomCommands                 []string
	InitDaggerWithWorkDirByDefault bool
	RunInVendor                    bool
}

func GetCLIGlobalArgs() CLIGlobalArgs {

	keValuePairsFromViper := viper.Get("set-env")
	var envKeyValuePairsToSet map[string]interface{}

	if keValuePairsFromViper == nil {
		envKeyValuePairsToSet = make(map[string]interface{})
	} else {
		envKeyValuePairsToSet = keValuePairsFromViper.(map[string]interface{})
	}

	args := CLIGlobalArgs{
		WorkingDir: viper.Get("work-dir").(string),
		MountDir:   viper.Get("mount-dir").(string),
		TargetDir:  viper.Get("target-dir").(string),
		TaskName:   viper.Get("task").(string),
		//ScanEnvVarKeys: viper.Get("scan-env").([]string),
		ScanEnvVarKeys: []string{},
		//EnvKeyValuePairsToSet:          viper.Get("set-env").(map[string]interface{}),
		EnvKeyValuePairsToSet: envKeyValuePairsToSet,
		ScanAWSKeys:           viper.Get("scan-aws-keys").(bool),
		ScanTerraformVars:     viper.Get("scan-terraform-vars").(bool),
		//CustomCommands:                 viper.Get("custom-cmds").([]string),
		CustomCommands:                 []string{},
		InitDaggerWithWorkDirByDefault: viper.Get("init-dagger-with-workdir").(bool),
		RunInVendor:                    viper.Get("run-in-vendor").(bool),
	}

	for k, v := range args.EnvKeyValuePairsToSet {
		args.EnvKeyValuePairsToSetString[k] = v.(string)
	}

	return args
}

func ShowCLITitle() {
	ux := tui.TUITitle{}
	ux.ShowTitleAndDescription("STILETTO",
		"Stiletto is a pipeline framework that works on top of Dagger.io. "+
			"Makes your pipelines more readable and easier to maintain.")
}
