package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type PageID string

type Page struct {
	ID           PageID  `json:"id"`
	Name         string  `json:"name"`
	Content      string  `json:"Content"`
	ParentPageID *PageID `json:"parent_page_id,omitempty"`
	Pages        []Page  `json:"pages"`
}

func (c *ClickupClient) GetDocPages(ctx context.Context, wks WorkspaceID, doc DocID) ([]Page, error) {
	if doc == "5fdw2-1681" {
		return []Page{}, nil
	}
	pages := []Page{}
	if err := c.request(ctx, "GET", fmt.Sprintf("/v3/workspaces/%s/docs/%s/pages", wks, doc), url.Values{}, nil, &pages); err != nil {
		var eu ErrUnexpectedStatusCode
		if errors.As(err, &eu) && eu == http.StatusNotFound {
			return pages, nil
		}

		return nil, fmt.Errorf("error getting document (doc=%s) pages: %w", doc, err)
	}

	return pages, nil
}
