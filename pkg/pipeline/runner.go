package pipeline

import (
	"context"
	"dagger.io/dagger"
	"github.com/Excoriate/stiletto/internal/logger"
	"github.com/Excoriate/stiletto/internal/tui"
	"github.com/Excoriate/stiletto/pkg/config"
)

type Config struct {
	Logger       logger.Logger
	Dirs         config.DefaultDirs
	UXDisplay    tui.TUIDisplayer
	UXMessage    tui.TUIMessenger
	Platforms    map[dagger.Platform]string
	PipelineOpts *config.PipelineOptions
	Ctx          context.Context
}
