package scan

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

// JSCrawler reads JS discovered on page and extracts wss:// urls
func JSCrawler(url string) error {
	var html string
	// super basic, can 1000% be improved.
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2000*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			rootNode, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			html, err = dom.GetOuterHTML().WithNodeID(rootNode.NodeID).Do(ctx)
			return err
		}),
	) // asså usch
	if err != nil {
		log.Fatal("Error automation logic: ", err)
	}

	fmt.Println(html)

	// Parse the HTML and extract my ws urls
	ws_proto_regex, err := regexp.Compile(`ws(s?):\/\/([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?`)
	if err != nil {
		fmt.Println("Regex search failed -> ", err)

	}
	fmt.Println(ws_proto_regex)
	return nil

}
