package models

const (
	ScrapeDataFormat = "02-01-2006"
)

// ScraperEvent represents an event when the scraper gets data
type ScraperEvent struct {
	SiteName string `json:"site_name"`
	Msg      string `json:"msg"`
	Date     string `json:"date"`
}
