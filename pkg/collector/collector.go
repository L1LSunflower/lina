package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"

	"github.com/L1LSunflower/lina/internal/entities"
)

const (
	attributeLink = "product"

	url          = "https://lichi.com"
	saleEndpoint = "/kz/ru/sale"

	getNamePath = "product-content_product_detail_heading__CLArW"

	staticUrl = "https://static.lichi.com/product/"
)

// TODO: scroller make scroll til the end of list
// TODO: optimize product scrapper
// TODO: get colors and color types

var tempUrls = []string{"https://lichi.com/kz/ru/product/46007", "https://lichi.com/kz/ru/product/46397", "https://lichi.com/kz/ru/product/46395", "https://lichi.com/kz/ru/product/46148"}

func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// Disable the headless mode to see what happen.
		chromedp.Flag("headless", false /*true*/),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx /*, chromedp.WithDebugf(log.Printf)*/)
	defer cancel()
	var nodes []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Navigate(url+saleEndpoint)); err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)
	if err := chromedp.Run(ctx, chromedp.Nodes("a", &nodes, chromedp.ByQueryAll)); err != nil {
		panic(err)
	}
	fmt.Println(nodes)
	var links []string
	for _, node := range nodes {
		if link := findInAttribute(node.Attributes, attributeLink); link != "" {
			links = append(links, url+link)
		}
	}
	tempUrls = removeDuplicates(tempUrls)
	fmt.Println("links length:", len(tempUrls))
	var items []*entities.Item
	for _, link := range tempUrls {
		item, err := getProductContent(ctx, link)
		if err != nil {
			fmt.Println("ID ITEM ERROR:", item, " ERROR:", err)
		}
		items = append(items, item)
	}
	//fmt.Println(items)
}

func getProductContent(ctx context.Context, url string) (*entities.Item, error) {
	if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
		return nil, err
	}
	var (
		name       string
		article    string
		nodes      []*cdp.Node
		childNodes []*cdp.Node
	)
	time.Sleep(2 * time.Second)
	// Get item name
	if err := chromedp.Run(ctx, chromedp.Text("h1", &name, chromedp.ByQuery)); err != nil {
		return nil, err
	}
	// Get item article
	if err := chromedp.Run(ctx, chromedp.Text("p", &article, chromedp.ByQuery)); err != nil {
		return nil, err
	}
	// Get price and currency
	if err := chromedp.Run(ctx, chromedp.Nodes(".product-content_product_sale_line__Cz1ea.ltr_mode.w-auto", &nodes, chromedp.ByQuery)); err != nil {
		return nil, err
	}
	if err := chromedp.Run(ctx, chromedp.Nodes("span", &childNodes, chromedp.ByQueryAll, chromedp.FromNode(nodes[0]))); err != nil {
		return nil, err
	}
	var priceStr, currency string
	if len(childNodes) >= 3 {
		priceStr = childNodes[1].Children[0].NodeValue
		currency = childNodes[2].Children[0].NodeValue
	}
	// Process price
	price, err := strToInt(priceStr)
	if err != nil {
		return nil, err
	}

	// Get actual price of item
	if err = chromedp.Run(ctx, chromedp.Nodes(".product-content_product_sale_price__xjnll", &nodes, chromedp.ByQuery)); err != nil {
		return nil, err
	}
	var actualPriceStr string
	if err = chromedp.Run(ctx, chromedp.Text("span", &actualPriceStr, chromedp.ByQuery, chromedp.FromNode(nodes[0]))); err != nil {
		return nil, err
	}
	// Process actual price
	actualPrice, err := strToInt(actualPriceStr)
	if err != nil {
		return nil, err
	}
	//// Get item sizes
	if err = chromedp.Run(ctx, chromedp.Nodes(".p-relative", &nodes, chromedp.ByQueryAll)); err != nil {
		return nil, err
	}
	var sizes []string
	for _, node := range nodes {
		if node.NodeName != "LI" {
			continue
		}
		size := ""
		if err = chromedp.Run(ctx, chromedp.Text("span", &size, chromedp.ByQuery, chromedp.FromNode(node))); err != nil {
			return nil, err
		}
		sizes = append(sizes, size)
	}
	// Get images
	if err = chromedp.Run(ctx, chromedp.Nodes(".ui-image_ui_image__ZWo6S", &nodes, chromedp.ByQueryAll)); err != nil {
		return nil, err
	}
	var imageLinks []string
	for _, node := range nodes {
		if link := findInAttribute(node.Attributes, getProductFromLink(url)); link != "" {
			imageLinks = append(imageLinks, link)
		}
	}
	for i := 0; i < len(imageLinks); i++ {
		if link := prepareLink(imageLinks[i]); link != "" {
			imageLinks[i] = link
		}
	}
	// Print out in stdout all grabbed data
	fmt.Printf("Наименование товара: %s\n", name)
	fmt.Printf("Артикул товара: %s\n", article)
	fmt.Printf("Цена без скидки: %d %s\n", price, currency)
	fmt.Printf("Цена с скидкой: %d %s\n", actualPrice, currency)
	fmt.Printf("Валюта: %s\n", currency)
	fmt.Printf("Размеры: %s\n", sizes)
	for _, link := range imageLinks {
		fmt.Printf("Ссылка на картинку: %s\n", link)
	}
	fmt.Println()
	//fmt.Printf("Цвет: %s\n\n", color)
	return &entities.Item{
		Url:           url,
		Name:          name,
		Article:       article,
		ExpectedPrice: price,
		ActualPrice:   actualPrice,
		Currency:      currency,
		Sizes:         sizes,
		ImageLinks:    imageLinks,
	}, nil
}

func getProductLink(attributes []string) string {
	for _, attribute := range attributes {
		if ind := strings.Index(attribute, attributeLink); ind >= 0 {
			return url + attribute
		}
	}
	return ""
}

func findInAttribute(attributes []string, pattern string) string {
	for _, attribute := range attributes {
		if i := strings.Index(attribute, pattern); i >= 0 {
			return attribute
		}
	}
	return ""
}

func strToInt(val string) (int, error) {
	number := ""
	for _, v := range val {
		if v >= '0' && v <= '9' {
			number += string(v)
		}
	}
	numb, err := strconv.Atoi(number)
	if err != nil {
		return 0, err
	}
	return numb, nil
}

func removeDuplicates(links []string) []string {
	linkMap := map[string]bool{}
	for _, link := range links {
		linkMap[link] = true
	}
	var result []string
	for link := range linkMap {
		result = append(result, link)
	}
	return result
}

func prepareLink(link string) string {
	if clearLink, _, ok := strings.Cut(link, "?"); ok {
		return clearLink
	}
	return ""
}

func getProductFromLink(link string) string {
	if s := strings.Split(link, "/"); len(s) > 2 {
		return "/" + s[len(s)-2] + "/" + s[len(s)-1]
	}
	return ""
}
