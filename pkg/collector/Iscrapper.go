package main

import "time"

type ScrapperInterface interface {
	InitInstance()
	Navigate(url string, timeDelay time.Duration) error
	Text(pattern string) (string, error)
	Nodes(pattern string)
}
