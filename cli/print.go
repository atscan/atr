package cli

import (
	"fmt"
	"io"
	"log"
	"os"

	jsoniter "github.com/json-iterator/go"
)

func Print(v interface{}) error {
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

func PrettyPrint(v interface{}, hg func(io.Writer, string)) error {
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
