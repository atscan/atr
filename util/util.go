package util

import "github.com/ipfs/go-cid"

func ReverseCidSlice(arr []cid.Cid) []cid.Cid {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}
