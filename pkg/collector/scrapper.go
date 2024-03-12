package main

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

// "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"

type Scrapper struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	Headless   bool
	DebugMode  bool
	UserAgent  string
}

func New(ctx context.Context, opts *ScrapperOptions) ScrapperInterface {
	return &Scrapper{
		ctx:       ctx,
		Headless:  opts.Headless,
		DebugMode: opts.DebugMode,
		UserAgent: opts.UserAgent,
	}
}

type ScrapperOptions struct {
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
	}
}

func (s *Scrapper) Navigate(url string, timeDelay time.Duration) error {
	if err := chromedp.Run(s.ctx, chromedp.Navigate(url)); err != nil {
		return err
	}
	time.Sleep(timeDelay)
	return nil
}

func (s *Scrapper) Text(pattern string) (string, error) {
	tempString := ""
	if err := chromedp.Run(s.ctx, chromedp.Text("h1", &tempString, chromedp.ByQuery)); err != nil {
		return "", err
	}
	return tempString, nil
}

func (s *Scrapper) Nodes(pattern string) {
	var tempNodes []*cdp.Node
	if err := chromedp.Run(s.ctx, chromedp.Nodes("a", &tempNodes, chromedp.ByQueryAll)); err != nil {
		panic(err)
	}
}
