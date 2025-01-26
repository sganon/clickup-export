package main

type WorkspaceID string

type Workspace struct {
	ID     WorkspaceID `json:"id"`
	Spaces []Space     `json:"-"`
	Docs   []Doc       `json:"-"`
}
