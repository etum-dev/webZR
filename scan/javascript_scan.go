package scan

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/etum-dev/WebZR/utils"
)

// JSCrawler reads JS discovered on page and extracts wss:// urls
func JSCrawler(domain string) []utils.ScanResult {
	//TODO: convert the URL so chromedp likes it
	pdomain := utils.AppendHttpsProto(domain)

	var html string

	// super basic, can 1000% be improved.
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(pdomain),
		chromedp.Sleep(2000*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			rootNode, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			html, err = dom.GetOuterHTML().WithNodeID(rootNode.NodeID).Do(ctx)
			return err
		}),
	)
	if err != nil {
		log.Print("Error automation logic: ", err, domain)

	}
	// TODO: search for stuff like WebSocket(url); and just report back that it potentially has one
	// Parse the HTML and extract my ws urls
	ws_proto_regex, err := regexp.Compile(`ws(s?):\/\/([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?`)
	if err != nil {
		fmt.Println("Regex search failed -> ", err)

	}

	fmt.Println(ws_proto_regex.FindString(html))
	extractedUrl := ws_proto_regex.FindString(html)

	var results []utils.ScanResult
	results = append(results, utils.ScanResult{
		StatusCode: 0,
		URL:        extractedUrl,
		Host:       domain,
		Scheme:     "javascript",
	})

	return results

}
