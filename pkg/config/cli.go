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
	ScanEnvVarsWithPrefix          []string
	DotEnvFile                     string
	ScanAllEnvVars                 bool
	CustomCommands                 []string
	InitDaggerWithWorkDirByDefault bool
	RunInVendor                    bool
}

func GetCLIGlobalArgs() (CLIGlobalArgs, error) {
	cfg := Cfg{}
	defaultEmptyMap := make(map[string]interface{})

	// 'set-env' option
	var setEnvValue map[string]interface{}
	keValuePairsFromViper, err := cfg.GetStringInterfaceMapFromViper("set-env")
	if err != nil {
		setEnvValue = defaultEmptyMap
	} else {
		setEnvValue = keValuePairsFromViper.Value.(map[string]interface{})
	}

	// 'scan-env' option
	var scanEnvVarKeys []string
	scanEnvVarKeysFromViper, err := cfg.GetStringSliceFromViper("scan-env")
	if err != nil {
		scanEnvVarKeys = []string{}
	} else {
		scanEnvVarKeys = scanEnvVarKeysFromViper.Value.([]string)
	}

	// Scan env vars with prefix
	var scanEnvVarsWithPrefix []string
	scanEnvVarsWithPrefixFromViper, err := cfg.GetStringSliceFromViper("scan-env-vars-prefix")
	if err != nil {
		scanEnvVarsWithPrefix = []string{}
	} else {
		scanEnvVarsWithPrefix = scanEnvVarsWithPrefixFromViper.Value.([]string)
	}

	args := CLIGlobalArgs{
		WorkingDir:            viper.GetString("work-dir"),
		MountDir:              viper.GetString("mount-dir"),
		TargetDir:             viper.GetString("target-dir"),
		TaskName:              viper.GetString("task"),
		ScanEnvVarKeys:        scanEnvVarKeys,
		EnvKeyValuePairsToSet: setEnvValue,
		ScanAWSKeys:           viper.GetBool("scan-aws-keys"),
		ScanTerraformVars:     viper.GetBool("scan-terraform-vars"),
		ScanEnvVarsWithPrefix: scanEnvVarsWithPrefix,
		ScanAllEnvVars:        viper.GetBool("scan-all-env-vars"),
		DotEnvFile:            viper.GetString("dot-env-file"),
		//CustomCommands:                 viper.Get("custom-cmds").([]string),
		CustomCommands:                 []string{},
		InitDaggerWithWorkDirByDefault: viper.GetBool("init-dagger-with-workdir"),
		RunInVendor:                    viper.GetBool("run-in-vendor"),
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
