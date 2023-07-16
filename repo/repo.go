package repo

import (
	"context"
	"fmt"
	"io"

	cbor2 "github.com/fxamacker/cbor/v2"

	"github.com/bluesky-social/indigo/mst"
	"github.com/bluesky-social/indigo/util"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/ipld/go-car/v2"
)

// current version of repo currently implemented
const ATP_REPO_VERSION int64 = 2

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

func IngestRepo(ctx context.Context, bs blockstore.Blockstore, r io.Reader) (cid.Cid, error) {
	br, err := car.NewBlockReader(r)
	if err != nil {
		return cid.Undef, err
	}

	for {
		blk, err := br.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return cid.Undef, err
		}

		if err := bs.Put(ctx, blk); err != nil {
			return cid.Undef, err
		}
	}

	return br.Roots[0], nil
}

func ReadRepoFromCar(ctx context.Context, r io.Reader) (*Repo, error) {
	bs := blockstore.NewBlockstore(datastore.NewMapDatastore())
	root, err := IngestRepo(ctx, bs, r)
	if err != nil {
		return nil, err
	}

	return OpenRepo(ctx, bs, root, false)
}

func OpenRepo(ctx context.Context, bs blockstore.Blockstore, root cid.Cid, fullRepo bool) (*Repo, error) {
	cst := util.CborStore(bs)

	var sc SignedCommit
	if err := cst.Get(ctx, root, &sc); err != nil {
		return nil, fmt.Errorf("loading root from blockstore: %w", err)
	}

	if sc.Version != ATP_REPO_VERSION {
		return nil, fmt.Errorf("unsupported repo version: %d", sc.Version)
	}

	return &Repo{
		sc:      sc,
		bs:      bs,
		cst:     cst,
		repoCid: root,
	}, nil
}

func (r *Repo) getMst(ctx context.Context) (*mst.MerkleSearchTree, error) {
	if r.mst != nil {
		return r.mst, nil
	}

	t := mst.LoadMST(r.cst, r.sc.Data)
	r.mst = t
	return t, nil
}

var ErrDoneIterating = fmt.Errorf("done iterating")

func (r *Repo) ForEach(ctx context.Context, prefix string, cb func(k string, v cid.Cid) error) error {
	t := mst.LoadMST(r.cst, r.sc.Data)

	if err := t.WalkLeavesFrom(ctx, prefix, cb); err != nil {
		if err != ErrDoneIterating {
			return err
		}
	}

	return nil
}

func (r *Repo) GetRecord(ctx context.Context, rpath string) (cid.Cid, interface{}, error) {
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

	//fmt.Println(hex.EncodeToString(blk.RawData()))

	var v interface{}
	cbor2.Unmarshal(blk.RawData(), &v)

	//fmt.Println(string(json))

	/*rec, err := lexutil.CborDecodeValue(blk.RawData())
	if err != nil {
		return cid.Undef, nil, err
	}*/

	return cc, v, nil
}
