package logger

import (
	"sort"

	"github.com/vineethkrishnan/dockit/pkg/models"
)

// Config configures the log engine thresholds.
type Config struct {
	WarningThresholdMB int64
}

var DefaultConfig = Config{
	WarningThresholdMB: 100, // Warn if a container's log file is over 100MB
}

// Engine processes log metrics
type Engine struct {
	Config Config
}

func NewEngine(cfg Config) *Engine {
	return &Engine{Config: cfg}
}

// AnalyzeLogSizes sorts the input and returns a warning flag for massive logs.
func (e *Engine) AnalyzeLogSizes(metrics []*models.LogMetrics) ([]*models.LogMetrics, int64) {
	var totalSize int64

	// Sort by Largest LogSize first
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].LogSize > metrics[j].LogSize
	})

	for _, m := range metrics {
		totalSize += m.LogSize
	}

	return metrics, totalSize
}

// IsExcessive converts the WarningThresholdMB to bytes and checks if the log exceeds it.
func (e *Engine) IsExcessive(sizeBytes int64) bool {
	return sizeBytes > (e.Config.WarningThresholdMB * 1024 * 1024)
}
