package daggerio

import (
	"dagger.io/dagger"
	"fmt"
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/Excoriate/stiletto/internal/errors"
)

type RegistryAuthOptions struct {
	RegistryAddress string
	RegistryUser    string
	RegistrySecret  DaggerSecret
}

func AuthWithRegistry(c *dagger.Client, container *dagger.Container,
	opt RegistryAuthOptions) (*dagger.Container,
	error) {

	userNormalised := common.NormaliseNoSpaces(opt.RegistryUser)
	addrNormalised := common.NormaliseNoSpaces(opt.RegistryAddress)

	if userNormalised == "" || addrNormalised == "" {
		return nil, errors.NewDaggerConfigurationError(fmt.Sprintf("Failed to auth with registry. "+
			"User: %s, addr: %s", userNormalised, addrNormalised), nil)
	}

	if container == nil {
		return nil, errors.NewDaggerConfigurationError(fmt.Sprintf("Failed to auth with registry. "+
			"Container is nil"), nil)
	}

	if c == nil {
		return nil, errors.NewDaggerConfigurationError(fmt.Sprintf("Failed to auth with registry. "+
			"Dagger client is nil"), nil)
	}

	registrySecret := c.SetSecret(opt.RegistrySecret.SecretId, opt.RegistrySecret.SecretValue)

	return c.Container().WithRegistryAuth(addrNormalised, userNormalised, registrySecret), nil
}
