package scraper

import "errors"

var (
	ErrPlayerNotFound = errors.New("player not found")
	ErrFetchURL    = errors.New("failed to fetch URL")
	ErrScrapeFailed   = errors.New("failed to scrape player data")
	ErrPositionNotFound = errors.New("player position not found")
)