package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// will handle different kind of urls supplied

func CheckDomain(domain string) string {
	fmt.Println("Checking: ", domain)
	/*if strings.HasSuffix(domain, ".*") {
		fmt.Println("Domain ends with .*, indicating wildcard TLD")

	} */
	if strings.HasPrefix(domain, "*") {
		// indicates we def want to do subdomain fuzz
		domain = strings.Trim(domain, "*.")
		fmt.Println("rem wildcard: ", domain)
		//fixedDomain, err := url.Parse(domain)
		/*if err != nil {
			fmt.Println("Invalid domain: ", domain)
		}*/
		return domain
	}
	// check if any other issues:
	/*fixedDomain, err := url.Parse(domain)
	if err != nil {
		fmt.Println("Invalid domain: ", domain)
	} */

	return domain
}

func AppendProto(inurl string) string {
	// no protocol, assume https
	if !strings.Contains(inurl, "://") {
		return "https://" + inurl
	}

	parsedUrl, err := url.Parse(inurl)
	if err != nil {
		fmt.Println("parse error:", err)
		return inurl
	}

	// Always use https for http/ws schemes
	if parsedUrl.Scheme == "ws" || parsedUrl.Scheme == "wss" {
		parsedUrl.Scheme = "https"
	}

	return parsedUrl.String()

}
