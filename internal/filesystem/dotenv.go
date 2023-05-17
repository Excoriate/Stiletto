package filesystem

import (
	"bufio"
	"fmt"
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/Excoriate/stiletto/internal/errors"
	"os"
	"strings"
)

func GetEnvVarsFromDotFile(filepath string) (map[string]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	env := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, errors.NewInternalPipelineError(fmt.Sprintf("Invalid line in ."+
				"env file: %s", line))
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = common.RemoveDoubleQuotes(value)

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(env) == 0 {
		return nil, errors.NewInternalPipelineError(fmt.Sprintf(" .env file %s is empty"))
	}

	return env, nil
}
