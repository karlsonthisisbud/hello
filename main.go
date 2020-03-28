package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gocolly/colly/v2"
)

type SKU struct {
	SKUID    string
	Prices   []Price
	Name     string
	ImageURL string
	SKUURL   string
}

type Price struct {
	Price    float64
	Unit     int
	UnitType string
	Currency string
}

func main() {
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
			p, _ := strconv.ParseFloat(trimFirstRune(v), 64)
			prices[i] = Price{
				Price:    p,
				Currency: "GBP",
			}
		}
		for i, v := range e.ChildTexts("span.pricing__type") {
			u := strings.Split(v[1:], " ")
			prices[i].Unit, _ = strconv.Atoi(u[1])
			prices[i].UnitType = u[0]
		}

		sku := SKU{
			SKUID:    e.ChildText(".product__sku")[5:],
			Name:     e.ChildText(".product__title"),
			Prices:   prices,
			ImageURL: e.ChildAttr(".img-fluid", "src"),
			SKUURL:   e.ChildAttr("a", "href"),
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

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
