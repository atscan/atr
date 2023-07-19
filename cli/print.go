package cli

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/atscan/atr/repo"
	jsoniter "github.com/json-iterator/go"
)

type ObjectOutput struct {
	Did   string      `json:"did"`
	Rkey  string      `json:"rkey"`
	Cid   string      `json:"cid"`
	Body  interface{} `json:"body"`
	Match interface{} `json:"-"`
}

func Print(v interface{}, ri repo.RepoItem, rs repo.RepoSnapshot) error {

	json, err := jsoniter.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}
	s := string(json)
	if s == "null" {
		return nil
	}
	fmt.Println(string(json))
	return nil
}

func PrettyPrint(v interface{}, ri repo.RepoItem, rs repo.RepoSnapshot, hg func(io.Writer, string)) error {

	json, err := jsoniter.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	s := string(json)
	if s == "null" {
		return nil
	}
	hg(os.Stdout, s+"\n")
	return nil
}
