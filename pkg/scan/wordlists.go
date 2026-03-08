package scan

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type wordlistCache struct {
	once    sync.Once
	entries []string
	err     error
	path    string
}

func (w *wordlistCache) load() ([]string, error) {
	w.once.Do(func() {
		file, err := os.Open(w.path)
		if err != nil {
			w.err = fmt.Errorf("cannot open %s: %w", w.path, err)
			return
		}
		defer file.Close()

		var entries []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			entry := strings.TrimSpace(scanner.Text())
			if entry == "" || strings.HasPrefix(entry, "#") {
				continue
			}
			entries = append(entries, entry)
		}

		if err := scanner.Err(); err != nil {
			w.err = err
			return
		}

		w.entries = entries
	})
	return w.entries, w.err
}

var (
	endpointCache  = wordlistCache{path: "ws-endpoints.txt"}
	subdomainCache = wordlistCache{path: "ws-subdomain.txt"}
)

func getEndpointWordlist() ([]string, error) {
	return endpointCache.load()
}

func getSubdomainWordlist() ([]string, error) {
	return subdomainCache.load()
}
