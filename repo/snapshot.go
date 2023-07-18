package repo

import (
	"context"
	"log"
	"strings"

	cid "github.com/ipfs/go-cid"
)

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
	Repo  Repo
}

func (ss *RepoSnapshot) GetCollectionStats(root string) (map[string]int, error) {
	rctx := context.TODO()

	rr := ss.Repo
	if root != "" && root != "b" {
		cid, err := cid.Parse(root)
		if err != nil {
			log.Fatalf("cannot parse CID: %s", root)
		}
		cr, err := OpenRepo(rctx, rr.BlockStore(), cid, 0)
		if err != nil {
			log.Fatal("cannot open repo")
		}
		rr = *cr
	}
	stats := make(map[string]int)
	if err := rr.ForEach(rctx, "", func(k string, v cid.Cid) error {
		col := strings.Split(k, "/")[0]
		_, ok := stats[col]
		if !ok {
			stats[col] = 0
		}
		stats[col]++
		return nil
	}); err != nil {
		return nil, err
	}
	return stats, nil
}

func (ss *RepoSnapshot) LoadItems(root string) error {
	rctx := context.TODO()

	rr := ss.Repo
	if root != "" && root != "b" {
		cid, err := cid.Parse(root)
		if err != nil {
			log.Fatalf("cannot parse CID: %s", root)
		}
		cr, err := OpenRepo(rctx, rr.BlockStore(), cid, 0)
		if err != nil {
			log.Fatal("cannot open repo")
		}
		rr = *cr
	}

	var out []RepoItem
	if err := rr.ForEach(rctx, "", func(k string, v cid.Cid) error {
		_, rec, err := rr.GetRecord(rctx, k)
		if err != nil {
			log.Println("Cannot get record:", v.String())
		}
		out = append(out, RepoItem{Cid: v, Path: k, Body: rec})
		ss.Items = out
		return nil
	}); err != nil {
		return err
	}
	return nil
}
