package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type SKU struct {
	SKUID        string
	Prices       []Price
	UnitQuantity int
	Name         string
	ImageURL     string
}

type Price struct {
	OrderQuantity int
	Discount      float64
	Price         float64
}

func main() {
	c := colly.NewCollector() // Get Collections List Page

	detailCollector := c.Clone() // Get Product Page
	d := c.Clone()               // Get Discounts

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.HasPrefix(link, "/collections") && HasOnePrefix(link) && !HasSortBy(link) && !strings.Contains(link, "/product-samples/") {
			e.Request.Visit(link)
		}

		if strings.HasPrefix(link, "/collections") && strings.Contains(link, "/products/") && !strings.Contains(link, "/product-samples/") {
			detailCollector.Visit(e.Request.AbsoluteURL(link))
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	d.OnRequest(func(r *colly.Request) {
		fmt.Println("d: Visiting", r.URL)
	})

	detailCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("DC Visiting", r.URL)
	})

	d.OnHTML("div", func(e *colly.HTMLElement) {
		ps := make([]Price, 0, 5)

		e.ForEach("table > tbody > tr", func(_ int, el *colly.HTMLElement) {
			row := el.ChildTexts("td")
			if len(row) == 3 {
				unitQuantity, _ := strconv.Atoi(row[0])
				discount, _ := strconv.ParseFloat(strings.ReplaceAll(row[1], "%", ""), 64)
				p, _ := strconv.ParseFloat(strings.ReplaceAll(row[2], "£", ""), 64)
				price := Price{
					Price:         p,
					OrderQuantity: unitQuantity,
					Discount:      discount,
				}
				ps = append(ps, price)
			}
		})
		e.Request.Ctx.Put("prices", ps)
	})

	detailCollector.OnHTML("body", func(e *colly.HTMLElement) {
		p, _ := strconv.ParseFloat(strings.ReplaceAll(e.ChildText("div.detail-price > span.price"), "£", ""), 64)
		price := Price{
			Price:         p,
			OrderQuantity: 1,
			Discount:      0.0,
		}

		UnitQuantity, _ := strconv.Atoi(e.ChildTexts("#variant-attributes-table > tbody > tr > td")[1])
		fmt.Println(e.ChildTexts("#variant-attributes-table > tbody > tr > td")[0])

		sku := SKU{
			SKUID:        e.ChildTexts("#variant-attributes-table > tbody > tr > td")[0],
			ImageURL:     e.ChildAttr("img.zoom-image", "src"),
			Prices:       []Price{price},
			UnitQuantity: UnitQuantity,
			Name:         e.ChildAttr("#product-info > div > h1", "content"),
		}

		groupName := e.ChildAttr("#product-bulk-discount-table-container", "data-bulkdiscountgroup")
		originalPrice := strings.ReplaceAll(e.ChildText("div.detail-price > span.price"), "£", "")

		nextUrl := "https://apps.cupsdirect.co.uk/bulkdiscountservice/BulkDiscountTable?group=" + groupName + "&price=%C2%A3" + originalPrice
		fmt.Println(nextUrl)

		ctx := colly.NewContext()
		ctx.Put("sku", sku)

		d.Request("GET", nextUrl, nil, ctx, nil)

		pp := ctx.GetAny("prices")
		if pp != nil {
			casted := pp.([]Price)
			sku.Prices = append(sku.Prices, casted...)
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")

		// Dump json to the standard output
		enc.Encode(sku)
	})

	c.Visit("https://cupsdirect.co.uk/")
}

func HasOnePrefix(URL string) bool {
	slashes := strings.Count(URL, "/")
	return slashes == 2
}

func HasSortBy(URL string) bool {
	sort_by := strings.Count(URL, "sort_by")
	return sort_by > 0
}
