package main

import (
	"context"
	"fmt"
	"net/url"
)

type FolderID string

type Folder struct {
	ID     FolderID `json:"id"`
	Name   string   `json:"folder"`
	Hidden bool     `json:"hidden"`
	Space  *Space   `json:"space,omitempty"`
	Docs   []*Doc   `json:"-"`
	Lists  []*List  `json:"-"`
}

func (c *ClickupClient) GetFolderByID(ctx context.Context, id FolderID) (*Folder, error) {
	var folder Folder
	if err := c.request(ctx, "GET", fmt.Sprintf("/v2/folder/%s", id), url.Values{}, nil, &folder); err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	return &folder, nil
}

func (c *ClickupClient) GetSpaceFolders(ctx context.Context, id SpaceID) ([]*Folder, error) {
	type responseBody struct {
		Folders []*Folder `json:"folders"`
	}
	var rb responseBody
	if err := c.request(ctx, "GET", fmt.Sprintf("/v2/space/%s/folder", id), url.Values{}, nil, &rb); err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	return rb.Folders, nil
}

func PopulateFolders(ctx context.Context, cup *ClickupClient, wks WorkspaceID, folder *Folder) error {
	var err error
	folder.Docs, err = cup.GetFolderDocs(ctx, wks, folder.ID)
	if err != nil {
		return fmt.Errorf("error getting folder docs: %w", err)
	}
	for _, d := range folder.Docs {
		d.Pages, err = cup.GetDocPages(ctx, wks, d.ID)
		if err != nil {
			return fmt.Errorf("error getting document page: %w", err)
		}
	}

	folder.Lists, err = cup.GetFolderLists(ctx, folder.ID)
	if err != nil {
		return fmt.Errorf("error getting folder lists: %w", err)
	}

	for _, l := range folder.Lists {
		if err := PopulateList(ctx, cup, wks, l); err != nil {
			return fmt.Errorf("error populating list: %w", err)
		}
	}

	return nil
}
