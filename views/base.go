package views

import "github.com/torbatti/sqlite-cache/models"

type HeadInfo struct {
	Title       string
	Description string
}

type PageInfo struct {
	// seo related
	HeadInfo HeadInfo

	Game models.Game
}
