package scan

import (
	"bufio"
	"errors"
	"net/http"
	"strings"
)

const (
	RemoteWordlist = "https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/refs/heads/main/data/domains.txt"
)

func GetRemote() error {
	resp, err := http.Get(RemoteWordlist)
	// or wildcard domain + subfinder integration?
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("bro is bannned from github???" + resp.Status)
	}
	defer resp.Body.Close()

	list := bufio.NewScanner(resp.Body)
	const tokenFuckery = 10 * 1024 * 1024 // mod as needed
	list.Buffer(make([]byte, 0, 64*1024), tokenFuckery)

	for list.Scan(){
		line := strings.TrimSpace(list.Text())
		if line == "" {
			continue
		}
	}

	return list.Err()

}

