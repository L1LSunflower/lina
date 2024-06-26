package scrapper

import (
	"time"

	"github.com/chromedp/cdproto/cdp"
)

type Driver interface {
	Close() error
	InitInstance()
	ProcessPagination() error
	CheckIfExists(className string) (bool, error)
	Navigate(url string, timeDelay time.Duration) error
	Text(pattern string, selector querySelector) (string, error)
	Nodes(pattern string, selector querySelector) ([]*cdp.Node, error)
	NodesFromNode(pattern string, selector querySelector, node *cdp.Node) ([]*cdp.Node, error)
	TextFromNode(pattern string, selector querySelector, node *cdp.Node) (string, error)
}
