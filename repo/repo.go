package repo

import (
	"context"
	"fmt"
	"io"

	cbor2 "github.com/fxamacker/cbor/v2"

	"github.com/bluesky-social/indigo/mst"
	"github.com/bluesky-social/indigo/util"
	blockstore "github.com/ipfs/boxo/blockstore"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/ipld/go-car/v2"
)

// current version of repo currently implemented
const ATP_REPO_VERSION int64 = 3

type SignedCommit struct {
	Did     string   `cborgen:"did"`
	Version int64    `cborgen:"version"`
	Prev    *cid.Cid `cborgen:"prev"`
	Data    cid.Cid  `cborgen:"data"`
	Sig     []byte   `cborgen:"sig"`
}

type UnsignedCommit struct {
	Did     string   `cborgen:"did"`
	Version int64    `cborgen:"version"`
	Prev    *cid.Cid `cborgen:"prev"`
	Data    cid.Cid  `cborgen:"data"`
}

type Repo struct {
	Blocks int

	sc  SignedCommit
	cst cbor.IpldStore
	bs  blockstore.Blockstore

	repoCid cid.Cid

	mst *mst.MerkleSearchTree

	dirty bool
}

// Returns a copy of commit without the Sig field. Helpful when verifying signature.
func (sc *SignedCommit) Unsigned() *UnsignedCommit {
	return &UnsignedCommit{
		Did:     sc.Did,
		Version: sc.Version,
		Prev:    sc.Prev,
		Data:    sc.Data,
	}
}

// returns bytes of the DAG-CBOR representation of object. This is what gets
// signed; the `go-did` library will take the SHA-256 of the bytes and sign
// that.
/*func (uc *UnsignedCommit) BytesForSigning() ([]byte, error) {
	buf := new(bytes.Buffer)
	/if err := uc.MarshalCBOR(buf); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}*/

func IngestRepo(ctx context.Context, bs blockstore.Blockstore, r io.Reader) (cid.Cid, int, error) {
	br, err := car.NewBlockReader(r)
	if err != nil {
		return cid.Undef, 0, err
	}

	size := 0
	for {
		blk, err := br.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return cid.Undef, size, err
		}

		if err := bs.Put(ctx, blk); err != nil {
			return cid.Undef, size, err
		}
		size++
	}

	return br.Roots[0], size, nil
}

func ReadRepoFromCar(ctx context.Context, r io.Reader) (*Repo, error) {
	bs := blockstore.NewBlockstore(datastore.NewMapDatastore())
	root, size, err := IngestRepo(ctx, bs, r)
	if err != nil {
		return nil, err
	}

	return OpenRepo(ctx, bs, root, size)
}

func OpenRepo(ctx context.Context, bs blockstore.Blockstore, root cid.Cid, size int) (*Repo, error) {
	cst := util.CborStore(bs)

	var sc SignedCommit
	if err := cst.Get(ctx, root, &sc); err != nil {
		return nil, fmt.Errorf("loading root from blockstore: %w", err)
	}

	if sc.Version > ATP_REPO_VERSION {
		return nil, fmt.Errorf("unsupported repo version: %d", sc.Version)
	}

	return &Repo{
		sc:      sc,
		bs:      bs,
		cst:     cst,
		repoCid: root,
		Blocks:  size,
	}, nil
}

func (r *Repo) GetCommitsPath(len int) ([]cid.Cid, error) {
	path := []cid.Cid{}
	path = append(path, r.repoCid)
	if r.sc.Prev != nil {
		getParentCommits(r, *r.sc.Prev, &path, len-1)
	}
	return path, nil
}

func getParentCommits(r *Repo, c cid.Cid, p *[]cid.Cid, len int) ([]cid.Cid, error) {
	var sc SignedCommit
	ctx := context.TODO()
	if err := r.cst.Get(ctx, c, &sc); err != nil {
		return nil, fmt.Errorf("loading root from blockstore: %w", err)
	}
	*p = append(*p, c)
	if len == 0 {
		return nil, nil
	} else {
		len--
	}
	if sc.Prev != nil {
		return getParentCommits(r, *sc.Prev, p, len)
	}
	return nil, nil
}

func (r *Repo) Head() cid.Cid {
	return r.repoCid
}

func (r *Repo) SignedCommit() SignedCommit {
	return r.sc
}

func (r *Repo) MerkleSearchTree() *mst.MerkleSearchTree {
	return r.mst
}

func (r *Repo) BlockStore() blockstore.Blockstore {
	return r.bs
}

func (r *Repo) getMst(ctx context.Context) (*mst.MerkleSearchTree, error) {
	if r.mst != nil {
		return r.mst, nil
	}

	t := mst.LoadMST(r.cst, r.sc.Data)
	r.mst = t
	return t, nil
}

func (r *Repo) MST() *mst.MerkleSearchTree {
	mst, _ := r.getMst(context.TODO())
	return mst
}

var ErrDoneIterating = fmt.Errorf("done iterating")

func (r *Repo) ForEach(ctx context.Context, prefix string, cb func(k string, v cid.Cid) error) error {
	t, _ := r.getMst(ctx)

	if err := t.WalkLeavesFrom(ctx, prefix, cb); err != nil {
		if err != ErrDoneIterating {
			return err
		}
	}

	return nil
}

func (r *Repo) GetRecord(ctx context.Context, rpath string) (cid.Cid, map[string]interface{}, error) {
	mst, err := r.getMst(ctx)
	if err != nil {
		return cid.Undef, nil, fmt.Errorf("getting repo mst: %w", err)
	}

	cc, err := mst.Get(ctx, rpath)
	if err != nil {
		return cid.Undef, nil, fmt.Errorf("resolving rpath within mst: %w", err)
	}

	blk, err := r.bs.Get(ctx, cc)
	if err != nil {
		return cid.Undef, nil, err
	}
	var v map[string]interface{}
	cbor2.Unmarshal(blk.RawData(), &v)

	return cc, v, nil
}
