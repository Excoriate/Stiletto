package common

import "fmt"

// ValidateTerragruntCommands validates if the provided commands are valid terragrunt commands
func ValidateTerragruntCommands(commands []string) error {
	validCommands := []string{"run-all", "plan", "apply", "destroy", "show", "init"}

	for _, command := range commands {
		if !IsValidCommand(command, validCommands) {
			return fmt.Errorf("'%s' is not a valid terragrunt command. Valid commands are: %v", command, validCommands)
		}
	}

	return nil
}

// IsValidCommand checks if a given command is in the list of valid commands
func IsValidCommand(command string, validCommands []string) bool {
	for _, validCommand := range validCommands {
		if command == validCommand {
			return true
		}
	}
	return false
}
