package stiletto

import (
	"context"
	"github.com/Excoriate/stiletto/internal/config"
	"github.com/Excoriate/stiletto/internal/logger"
)

type Stiletto struct {
	Context context.Context
	Logger  logger.Logger
	Dirs    *config.DefaultDirs
}

func Init() Stiletto {
	log := &logger.StilettoLog{}
	log.InitLogger()

	s := Stiletto{
		Logger:  log,
		Dirs:    config.GetDefaultDirs(),
		Context: context.Background(),
	}

	return s
}
