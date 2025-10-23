package utils

import (
	"fmt"
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
