package lichi

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/L1LSunflower/lina/internal/drivers"
	"github.com/L1LSunflower/lina/internal/entities"

	"github.com/L1LSunflower/lina/pkg/scrapper"
)

const (
	defaultTimeout = 2 * time.Second

	attributeLink      = "product"
	linksPattern       = "a"
	h1Pattern          = "h1"
	pPattern           = "p"
	pricePattern       = ".product-content_product_sale_line__Cz1ea.ltr_mode.w-auto"
	spanPattern        = "span"
	actualPricePattern = ".product-content_product_sale_price__xjnll"
	pRelativePattern   = ".p-relative"
	imagePattern       = ".ui-image_ui_image__ZWo6S"

	url          = "https://lichi.com"
	saleEndpoint = "/kz/ru/sale"

	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"
)

func Items(ctx context.Context) ([]*entities.Item, error) {
	// Create new scrapper instance
	sc := scrapper.New(ctx, &scrapper.Options{Headless: true, DebugMode: false, UserAgent: userAgent})
	// Close with defer
	defer sc.Close()
	// Init chrome instance with options
	sc.InitInstance()
	products, err := allProducts(sc)
	if err != nil {
		return nil, err
	}
	var (
		items   []*entities.Item
		newItem *entities.Item
	)
	for _, productLink := range products {
		newItem, err = item(sc, productLink)
		if err != nil {
			log.Println("TIME: ", time.Now().Format(time.DateTime), " | ERROR: Get new item with error:", err)
			continue
		}
		items = append(items, newItem)
	}
	return items, nil
}

func allProducts(sc scrapper.ScrapperInterface) ([]string, error) {
	// Navigate to website
	if err := sc.Navigate(url+saleEndpoint, defaultTimeout); err != nil {
		return nil, err
	}
	// Wait til all products will be visible
	if err := sc.ProcessPagination(); err != nil {
		return nil, err
	}
	// Get all links on products from website
	nodes, err := sc.Nodes(linksPattern, scrapper.ByQueryAll)
	if err != nil {
		return nil, err
	}
	// Process links
	var links []string
	// Find links
	for _, node := range nodes {
		if link := drivers.FindInAttribute(node.Attributes, attributeLink); link != "" {
			links = append(links, url+link)
		}
	}
	// Remove duplicated links
	links = drivers.RemoveDuplicates(links)
	return links, nil
}

func item(sc scrapper.ScrapperInterface, link string) (*entities.Item, error) {
	var (
		newItem = &entities.Item{}
		err     error
	)
	// Navigate to website
	if err = sc.Navigate(link, defaultTimeout); err != nil {
		return nil, err
	}
	// Get item name
	if newItem.Name, err = sc.Text(h1Pattern, scrapper.ByQuery); err != nil {
		return nil, err
	}
	// Get item article
	if newItem.Article, err = sc.Text(pPattern, scrapper.ByQuery); err != nil {
		return nil, err
	}
	// Get price and currency
	if err = priceAndCurrency(sc, newItem); err != nil {
		return nil, err
	}
	// Get actual price
	if err = actualPrice(sc, newItem); err != nil {
		return nil, err
	}
	// Get sizes
	if newItem.Sizes, err = sizes(sc); err != nil {
		return nil, err
	}
	// Get image links
	if newItem.ImageLinks, err = imageLinks(sc); err != nil {
		return nil, err
	}
	return newItem, nil
}

func priceAndCurrency(sc scrapper.ScrapperInterface, newItem *entities.Item) error {
	// Create error variable
	var err error
	// If panic recover it with defer
	defer func(err error) {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered panic: %v", r)
		}
	}(err)
	// Get nodes with price
	nodes, err := sc.Nodes(pricePattern, scrapper.ByQuery)
	if err != nil {
		return err
	}
	// Get price and currency string
	nodes, err = sc.NodesFromNode(spanPattern, scrapper.ByQueryAll, nodes[0])
	if err != nil {
		return err
	}
	// Check nodes length
	if len(nodes) < 3 {
		return fmt.Errorf("nodes length less than needed (3)")
	}
	// Convert string price to int
	if newItem.ExpectedPrice, err = drivers.StrToInt(nodes[1].Children[0].NodeValue); err != nil {
		return err
	}
	// Currency from node as string
	newItem.Currency = nodes[2].Children[0].NodeValue
	// Return error if panic
	return err
}

func actualPrice(sc scrapper.ScrapperInterface, newItem *entities.Item) error {
	// Get actual price node
	nodes, err := sc.Nodes(actualPricePattern, scrapper.ByQuery)
	if err != nil {
		return err
	}
	// Check nodes length
	if len(nodes) < 1 {
		return fmt.Errorf("nodes length less than needed (1)")
	}
	// Get price from node
	priceStr, err := sc.TextFromNode(spanPattern, scrapper.ByQuery, nodes[0])
	if err != nil {
		return err
	}
	// Convert price string to int
	if newItem.ActualPrice, err = drivers.StrToInt(priceStr); err != nil {
		return err
	}
	return nil
}

func sizes(sc scrapper.ScrapperInterface) ([]string, error) {
	// Get item sizes
	nodes, err := sc.Nodes(pRelativePattern, scrapper.ByQueryAll)
	if err != nil {
		return nil, err
	}
	// Declare variables
	var (
		itemSizes []string
		size      string
	)
	// Process nodes to get sizes
	for _, node := range nodes {
		// Check if node name not "LI" skip that node
		if node.NodeName != "LI" {
			continue
		}
		// Get link text from node
		if size, err = sc.TextFromNode(spanPattern, scrapper.ByQuery, node); err != nil {
			itemSizes = append(itemSizes, size)
		}
	}
	// Return item sizes
	return itemSizes, nil
}

func imageLinks(sc scrapper.ScrapperInterface) ([]string, error) {
	// Create error variable
	var err error
	// If panic recover it with defer
	defer func(err error) {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered panic: %v", r)
		}
	}(err)
	// Get image nodes
	nodes, err := sc.Nodes(imagePattern, scrapper.ByQueryAll)
	if err != nil {
		return nil, err
	}
	var itemImageLinks []string
	for _, node := range nodes {
		if link := drivers.FindInAttribute(node.Attributes, drivers.ProductFromLink(url)); link != "" {
			itemImageLinks = append(itemImageLinks, link)
		}
	}
	for i := 0; i < len(itemImageLinks); i++ {
		if link := drivers.PrepareLink(itemImageLinks[i]); link != "" {
			itemImageLinks[i] = link
		}
	}
	return itemImageLinks, nil
}
