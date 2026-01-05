package scan

import (
	"github.com/etum-dev/WebZR/utils"
)

type Scanner interface {
	Scan(domain string) []utils.ScanResult
	Name() string
}
