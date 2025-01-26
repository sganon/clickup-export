package main

import (
	"context"
	"fmt"
	"net/url"
)

type DocID string

type Doc struct {
	ID       DocID        `json:"id"`
	Name     string       `json:"name"`
	Parent   ParentObject `json:"parent"`
	Archived bool         `json:"archived"`
	Deleted  bool         `json:"deleted"`
	Type     int          `json:"type"`
	Creator  int          `json:"creator"`
	Pages    []Page       `json:"-"`
}

func (c *ClickupClient) GetSpaceDocs(ctx context.Context, wks WorkspaceID, spaceID SpaceID) ([]*Doc, error) {
	return c.SearchDocs(ctx, wks, url.Values{
		"parent_type": []string{"SPACE"},
		"parent_id":   []string{string(spaceID)},
	})
}

func (c *ClickupClient) GetListDocs(ctx context.Context, wks WorkspaceID, listID ListID) ([]*Doc, error) {
	return c.SearchDocs(ctx, wks, url.Values{
		"parent_type": []string{"LIST"},
		"parent_id":   []string{string(listID)},
	})
}

func (c *ClickupClient) GetFolderDocs(ctx context.Context, wks WorkspaceID, folderID FolderID) ([]*Doc, error) {
	return c.SearchDocs(ctx, wks, url.Values{
		"parent_type": []string{"FOLDER"},
		"parent_id":   []string{string(folderID)},
	})
}

func (c *ClickupClient) SearchDocs(ctx context.Context, wks WorkspaceID, params url.Values) ([]*Doc, error) {
	docs := []*Doc{}
	if err := c.recursiveSearchDoc(ctx, wks, params, &docs, ""); err != nil {
		return nil, fmt.Errorf("error making a recursive search: %w", err)
	}

	return docs, nil
}

func (c *ClickupClient) recursiveSearchDoc(ctx context.Context, wks WorkspaceID, params url.Values, docs *[]*Doc, cursor string) error {
	type responseBody struct {
		Docs       []*Doc `json:"docs"`
		NextCursor string `json:"next_cursor"`
	}
	var rb responseBody
	params.Set("next_cursor", cursor)
	if err := c.request(ctx, "GET", fmt.Sprintf("/v3/workspaces/%s/docs", wks), params, nil, &rb); err != nil {
		return fmt.Errorf("error searching docs: %w", err)
	}
	*docs = append(*docs, rb.Docs...)
	if rb.NextCursor != "" {
		return c.recursiveSearchDoc(ctx, wks, params, docs, rb.NextCursor)
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
