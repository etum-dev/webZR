package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// will handle different kind of urls supplied

func CheckDomain(domain string) string {
	domain = strings.TrimSpace(domain)
	fmt.Println("Checking:", domain)

	if domain == "" {
		return domain
	}

	if strings.HasPrefix(domain, "*.") {
		domain = strings.TrimPrefix(domain, "*.")
		fmt.Println("rem wildcard:", domain)
	} else if strings.HasPrefix(domain, "*") {
		domain = strings.TrimPrefix(domain, "*")
		domain = strings.TrimPrefix(domain, ".")
		fmt.Println("rem wildcard:", domain)
	}

	domain = strings.TrimPrefix(domain, "//")

	if strings.Contains(domain, "://") {
		if parsed, err := url.Parse(domain); err == nil {
			domain = rebuildDomain(parsed)
		} else if idx := strings.Index(domain, "://"); idx != -1 {
			domain = domain[idx+3:]
		}
	}

	domain = strings.TrimLeft(domain, "/")
	return domain
}

func rebuildDomain(u *url.URL) string {
	host := u.Host
	if host == "" {
		host = u.Path
	}

	var builder strings.Builder
	builder.WriteString(strings.Trim(host, "/"))

	path := strings.Trim(u.Path, "/")
	if host != "" {
		path = strings.TrimPrefix(path, host)
	}
	if path != "" {
		if !strings.HasPrefix(path, "/") {
			builder.WriteString("/")
		}
		builder.WriteString(path)
	}

	if u.RawQuery != "" {
		builder.WriteString("?")
		builder.WriteString(u.RawQuery)
	}

	if u.Fragment != "" {
		builder.WriteString("#")
		builder.WriteString(u.Fragment)
	}

	return strings.Trim(builder.String(), "/")
}

func AppendProto(inurl string) string {
	inurl = strings.TrimSpace(inurl)
	if inurl == "" {
		return ""
	}

	inurl = strings.TrimPrefix(inurl, "//")

	// assuming protocol here, unsure which is best
	if !strings.Contains(inurl, "://") {
		inurl = "wss://" + inurl
	}

	parsedUrl, err := url.Parse(inurl)
	if err != nil {
		fmt.Println("parse error:", err)
		return inurl
	}

	// Always use https for http/ws schemes
	switch parsedUrl.Scheme {
	case "ws", "wss", "http":
		parsedUrl.Scheme = "https"
	}

	return parsedUrl.String()
}

// ExtractHostname returns only the host:port portion of a domain or URL.
func ExtractHostname(domain string) string {
	withProto := AppendProto(domain)
	if withProto == "" {
		return ""
	}

	parsed, err := url.Parse(withProto)
	if err != nil {
		parts := strings.Split(domain, "/")
		return strings.TrimSpace(parts[0])
	}

	if parsed.Host != "" {
		return parsed.Host
	}

	return strings.Split(strings.Trim(parsed.Path, "/"), "/")[0]
}
