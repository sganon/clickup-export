package main

import (
	"context"
	"fmt"
	"net/url"
)

type SpaceID string

type Space struct {
	ID      SpaceID   `json:"id"`
	Name    string    `json:"name"`
	Folders []*Folder `json:"-"`
	Lists   []*List   `json:"-"`
	Docs    []*Doc    `json:"-"`
}

func (c *ClickupClient) GetWorkspaceSpaces(ctx context.Context, id WorkspaceID) ([]*Space, error) {
	type responseBody struct {
		Spaces []*Space `json:"spaces"`
	}
	var rb responseBody
	if err := c.request(ctx, "GET", fmt.Sprintf("/v2/team/%s/space", id), url.Values{}, nil, &rb); err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	return rb.Spaces, nil
}

func (c *ClickupClient) GetSpaceByID(ctx context.Context, id SpaceID) (*Space, error) {
	var space Space
	if err := c.request(ctx, "GET", fmt.Sprintf("/v2/space/%s", id), url.Values{}, nil, &space); err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	return &space, nil
}

func PopulateSpace(ctx context.Context, cup *ClickupClient, wks WorkspaceID, space *Space) error {
	var err error
	space.Docs, err = cup.GetSpaceDocs(ctx, wks, space.ID)
	if err != nil {
		return fmt.Errorf("error getting docs: %w", err)
	}
	for _, d := range space.Docs {
		d.Pages, err = cup.GetDocPages(ctx, wks, d.ID)
		if err != nil {
			return fmt.Errorf("error getting document page: %w", err)
		}
	}

	space.Lists, err = cup.GetSpaceLists(ctx, space.ID)
	if err != nil {
		return fmt.Errorf("error getting lists: %w", err)
	}
	for _, l := range space.Lists {
		if err := PopulateList(ctx, cup, wks, l); err != nil {
			return fmt.Errorf("error populating list: %w", err)
		}
	}

	space.Folders, err = cup.GetSpaceFolders(ctx, space.ID)
	if err != nil {
		return fmt.Errorf("error getting folders: %w", err)
	}
	for _, f := range space.Folders {
		if err := PopulateFolders(ctx, cup, wks, f); err != nil {
			return fmt.Errorf("error populating folders: %w", err)
		}
	}

	return nil
}
