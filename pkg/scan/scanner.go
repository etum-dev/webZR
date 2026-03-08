package scan

import (
	"github.com/etum-dev/WebZR/pkg/utils"
)

type Scanner interface {
	Scan(domain string) []utils.ScanResult
	Name() string
}
