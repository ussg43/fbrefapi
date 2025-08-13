package scraper

import "errors"

var (
	ErrPlayerNotFound = errors.New("player not found")
	ErrInvalidClub    = errors.New("invalid club")
	ErrInvalidName    = errors.New("invalid player name")
	ErrScrapeFailed   = errors.New("failed to scrape player data")
	ErrP90NotAvailable = errors.New("P90 stats not available for this player")
	ErrSeasonNotFound = errors.New("season not found")
)