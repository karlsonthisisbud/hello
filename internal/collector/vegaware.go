package collector

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/i25959341/sku-aggregator/internal/types"
)

func CollectVegaware() []types.SKU {
	skus := make([]SKU, 0, 200)

	c := colly.NewCollector()

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// If attribute class is this long string return from callback
		// As this a is irrelevant

		if e.Attr("class") != "catLeft__single" {
			return
		}

		link := e.Attr("href")
		// If link start with browse or includes either signup or login return from callback
		if !strings.HasPrefix(link, "/uk/catalogue/") {
			return
		}
		// start scaping the page under the link found
		e.Request.Visit(link)
	})

	c.OnHTML(".product", func(e *colly.HTMLElement) {

		prices := make([]Price, len(e.ChildTexts("span.pricing__price")))

		for i, v := range e.ChildTexts("span.pricing__price") {
			prices[i] = Price{
				Price: v,
			}
		}
		fmt.Println(e.ChildTexts("span.pricing__type")[0])
		for i, v := range e.ChildTexts("span.pricing__type") {
			prices[i].Unit = v
		}

		sku := SKU{
			SKUID:  e.ChildText(".product__sku"),
			Name:   e.ChildText(".product__title"),
			Prices: prices,
		}
		skus = append(skus, sku)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.vegware.com/uk/catalogue/double_wall_cups/")

	jsonData, err := json.MarshalIndent(skus, "", "  ")
	if err != nil {
		panic(err)
	}

	// Dump json to the standard output (can be redirected to a file)
	fmt.Println(string(jsonData))
}
