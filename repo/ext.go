package repo

import cid "github.com/ipfs/go-cid"

type RepoItem struct {
	Cid  cid.Cid
	Path string
	Body interface{}
}

type RepoSnapshot struct {
	Root  cid.Cid
	File  string
	Size  int
	Items []RepoItem
}
