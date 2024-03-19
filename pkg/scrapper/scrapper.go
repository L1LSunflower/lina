package scrapper

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

type querySelector string

func (q querySelector) Selector() func(s *chromedp.Selector) {
	switch q {
	case ByQuery:
		return chromedp.ByQuery
	case ByQueryAll:
		return chromedp.ByQueryAll
	case ByJsPath:
		return chromedp.ByJSPath
	default:
		return chromedp.ByQuery
	}
}

const (
	ByQuery    querySelector = "byQuery"
	ByQueryAll querySelector = "byQueryAll"
	ByJsPath   querySelector = "byJsPath"
)

type Scrapper struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	Headless   bool
	DebugMode  bool
	UserAgent  string
}

func New(ctx context.Context, opts *Options) Driver {
	return &Scrapper{
		ctx:       ctx,
		Headless:  opts.Headless,
		DebugMode: opts.DebugMode,
		UserAgent: opts.UserAgent,
	}
}

type Options struct {
	Headless  bool
	DebugMode bool
	UserAgent string
}

func (s *Scrapper) InitInstance() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", s.Headless),
		chromedp.UserAgent(s.UserAgent),
	)
	s.ctx, s.cancelFunc = chromedp.NewExecAllocator(context.Background(), opts...)
	if s.DebugMode {
		s.ctx, s.cancelFunc = chromedp.NewContext(s.ctx, chromedp.WithDebugf(log.Printf))
	} else {
		s.ctx, s.cancelFunc = chromedp.NewContext(s.ctx)
	}
}

func (s *Scrapper) Navigate(url string, timeDelay time.Duration) error {
	if err := chromedp.Run(s.ctx, chromedp.Navigate(url)); err != nil {
		return err
	}
	time.Sleep(timeDelay)
	return nil
}

func (s *Scrapper) ProcessPagination() error {
	return nil
}

func (s *Scrapper) CheckIfExists(className string) (bool, error) {
	var exist bool
	if err := chromedp.Run(s.ctx, chromedp.EvaluateAsDevTools("document.getElementsByClassName('"+className+"').length > 0", &exist)); err != nil {
		return exist, err
	}
	return exist, nil
}

func (s *Scrapper) Nodes(pattern string, selector querySelector) ([]*cdp.Node, error) {
	var tempNodes []*cdp.Node
	if err := chromedp.Run(s.ctx, chromedp.Nodes(pattern, &tempNodes, selector.Selector())); err != nil {
		return nil, err
	}
	return tempNodes, nil
}

func (s *Scrapper) Text(pattern string, selector querySelector) (string, error) {
	tempString := ""
	if err := chromedp.Run(s.ctx, chromedp.Text(pattern, &tempString, selector.Selector())); err != nil {
		return "", err
	}
	return tempString, nil
}

func (s *Scrapper) NodesFromNode(pattern string, selector querySelector, node *cdp.Node) ([]*cdp.Node, error) {
	var tempNodes []*cdp.Node
	if err := chromedp.Run(s.ctx, chromedp.Nodes(pattern, &tempNodes, chromedp.FromNode(node), selector.Selector())); err != nil {
		return nil, err
	}
	return tempNodes, nil
}

func (s *Scrapper) TextFromNode(pattern string, selector querySelector, node *cdp.Node) (string, error) {
	tempString := ""
	if err := chromedp.Run(s.ctx, chromedp.Text(pattern, &tempString, selector.Selector(), chromedp.FromNode(node))); err != nil {
		return "", err
	}
	return tempString, nil
}

func (s *Scrapper) Close() error {
	return chromedp.Cancel(s.ctx)
}

/*
chromedp.ActionFunc(func(ctx context.Context) error {
    _, exp, err := runtime.Evaluate(`window.scrollTo(0,document.body.scrollHeight);`).Do(ctx)
    if err != nil {
        return err
    }
    if exp != nil {
        return exp
    }
    return nil

})
*/

//cdp.EvaluateAsDevTools("document.getElementsByClassName('product-content_product_sale_line__Cz1ea ltr_mode w-auto').length > 0", value)
