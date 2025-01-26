package main

import (
	"context"
	"fmt"
	"net/url"
)

type ListID string

type List struct {
	ID     ListID  `json:"id"`
	Name   string  `json:"name"`
	Folder *Folder `json:"folder"`
	Space  *Space  `json:"space"`
	Docs   []*Doc  `json:"-"`
}

func (c *ClickupClient) GetListByID(ctx context.Context, id ListID) (*List, error) {
	var list List
	if err := c.request(ctx, "GET", fmt.Sprintf("/v2/list/%s", id), url.Values{}, nil, &list); err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	return &list, nil
}

type searchListResponse struct {
	Lists []*List `json:"lists"`
}

func (c *ClickupClient) GetFolderLists(ctx context.Context, folderID FolderID) ([]*List, error) {
	var rb searchListResponse
	if err := c.request(ctx, "GET", fmt.Sprintf("/v2/folder/%s/list", folderID), url.Values{}, nil, &rb); err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	return rb.Lists, nil
}

func (c *ClickupClient) GetSpaceLists(ctx context.Context, spaceID SpaceID) ([]*List, error) {
	var rb searchListResponse
	if err := c.request(ctx, "GET", fmt.Sprintf("/v2/space/%s/list", spaceID), url.Values{}, nil, &rb); err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	return rb.Lists, nil
}

func PopulateList(ctx context.Context, cup *ClickupClient, wks WorkspaceID, list *List) error {
	var err error
	list.Docs, err = cup.GetListDocs(ctx, wks, list.ID)
	if err != nil {
		return fmt.Errorf("error getting docs: %w", err)
	}
	for _, d := range list.Docs {
		d.Pages, err = cup.GetDocPages(ctx, wks, d.ID)
		if err != nil {
			return fmt.Errorf("error getting document page: %w", err)
		}
	}
	return nil
}
