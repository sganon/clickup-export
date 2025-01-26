package main

import (
	"context"
	"fmt"
	"net/url"
)

type DocID string

type Doc struct {
	ID     DocID        `json:"id"`
	Name   string       `json:"name"`
	Parent ParentObject `json:"parent"`
}

func (c *ClickupClient) SearchDocs(ctx context.Context, wks WorkspaceID) ([]Doc, error) {
	docs := []Doc{}
	if err := c.recursiveSearchDoc(ctx, wks, &docs, ""); err != nil {
		return nil, fmt.Errorf("error making a recursive search: %w", err)
	}

	return docs, nil
}

func (c *ClickupClient) recursiveSearchDoc(ctx context.Context, wks WorkspaceID, docs *[]Doc, cursor string) error {
	type responseBody struct {
		Docs       []Doc  `json:"docs"`
		NextCursor string `json:"next_cursor"`
	}
	var rb responseBody
	if err := c.request(ctx, "GET", fmt.Sprintf("/workspaces/%s/docs", wks), url.Values{
		"next_cursor": []string{cursor},
	}, nil, &rb); err != nil {
		return fmt.Errorf("error searching docs: %w", err)
	}
	*docs = append(*docs, rb.Docs...)
	if rb.NextCursor != "" {
		return c.recursiveSearchDoc(ctx, wks, docs, rb.NextCursor)
	}

	return nil
}

type ParentObjectType = int

const (
	ParentSpace      ParentObjectType = 4
	ParentFolder     ParentObjectType = 5
	ParentList       ParentObjectType = 6
	ParentEverything ParentObjectType = 7
	ParentWorkspace  ParentObjectType = 12
)

type ParentObject struct {
	ID   string           `json:"id"`
	Type ParentObjectType `json:"type"`
}
