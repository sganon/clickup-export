package main

type PageID string

type Page struct {
	ID           PageID  `json:"id"`
	Name         string  `json:"name"`
	Content      string  `json:"Content"`
	ParentPageID *PageID `json:"parent_page_id,omitempty"`
	Pages        []Page  `json:"pages"`
}
