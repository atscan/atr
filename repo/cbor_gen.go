// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package repo

import (
	"fmt"
	"io"
	"math"
	"sort"

	cid "github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf
var _ = cid.Undef
var _ = math.E
var _ = sort.Sort

func (t *SignedCommit) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{165}); err != nil {
		return err
	}

	// t.Did (string) (string)
	if len("did") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"did\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("did"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("did")); err != nil {
		return err
	}

	if len(t.Did) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Did was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Did))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Did)); err != nil {
		return err
	}

	// t.Sig ([]uint8) (slice)
	if len("sig") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"sig\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("sig"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("sig")); err != nil {
		return err
	}

	if len(t.Sig) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.Sig was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.Sig))); err != nil {
		return err
	}

	if _, err := cw.Write(t.Sig[:]); err != nil {
		return err
	}

	// t.Data (cid.Cid) (struct)
	if len("data") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"data\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("data"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("data")); err != nil {
		return err
	}

	if err := cbg.WriteCid(cw, t.Data); err != nil {
		return xerrors.Errorf("failed to write cid field t.Data: %w", err)
	}

	// t.Prev (cid.Cid) (struct)
	if len("prev") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"prev\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("prev"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("prev")); err != nil {
		return err
	}

	if t.Prev == nil {
		if _, err := cw.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCid(cw, *t.Prev); err != nil {
			return xerrors.Errorf("failed to write cid field t.Prev: %w", err)
		}
	}

	// t.Version (int64) (int64)
	if len("version") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"version\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("version"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("version")); err != nil {
		return err
	}

	if t.Version >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Version)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.Version-1)); err != nil {
			return err
		}
	}
	return nil
}

func (t *SignedCommit) UnmarshalCBOR(r io.Reader) (err error) {
	*t = SignedCommit{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("SignedCommit: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadString(cr)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.Did (string) (string)
		case "did":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Did = string(sval)
			}
			// t.Sig ([]uint8) (slice)
		case "sig":

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}

			if extra > cbg.ByteArrayMaxLen {
				return fmt.Errorf("t.Sig: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra > 0 {
				t.Sig = make([]uint8, extra)
			}

			if _, err := io.ReadFull(cr, t.Sig[:]); err != nil {
				return err
			}
			// t.Data (cid.Cid) (struct)
		case "data":

			{

				c, err := cbg.ReadCid(cr)
				if err != nil {
					return xerrors.Errorf("failed to read cid field t.Data: %w", err)
				}

				t.Data = c

			}
			// t.Prev (cid.Cid) (struct)
		case "prev":

			{

				b, err := cr.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := cr.UnreadByte(); err != nil {
						return err
					}

					c, err := cbg.ReadCid(cr)
					if err != nil {
						return xerrors.Errorf("failed to read cid field t.Prev: %w", err)
					}

					t.Prev = &c
				}

			}
			// t.Version (int64) (int64)
		case "version":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.Version = int64(extraI)
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
func (t *UnsignedCommit) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{164}); err != nil {
		return err
	}

	// t.Did (string) (string)
	if len("did") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"did\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("did"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("did")); err != nil {
		return err
	}

	if len(t.Did) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Did was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Did))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Did)); err != nil {
		return err
	}

	// t.Data (cid.Cid) (struct)
	if len("data") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"data\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("data"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("data")); err != nil {
		return err
	}

	if err := cbg.WriteCid(cw, t.Data); err != nil {
		return xerrors.Errorf("failed to write cid field t.Data: %w", err)
	}

	// t.Prev (cid.Cid) (struct)
	if len("prev") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"prev\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("prev"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("prev")); err != nil {
		return err
	}

	if t.Prev == nil {
		if _, err := cw.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCid(cw, *t.Prev); err != nil {
			return xerrors.Errorf("failed to write cid field t.Prev: %w", err)
		}
	}

	// t.Version (int64) (int64)
	if len("version") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"version\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("version"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("version")); err != nil {
		return err
	}

	if t.Version >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Version)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.Version-1)); err != nil {
			return err
		}
	}
	return nil
}

func (t *UnsignedCommit) UnmarshalCBOR(r io.Reader) (err error) {
	*t = UnsignedCommit{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("UnsignedCommit: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadString(cr)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.Did (string) (string)
		case "did":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Did = string(sval)
			}
			// t.Data (cid.Cid) (struct)
		case "data":

			{

				c, err := cbg.ReadCid(cr)
				if err != nil {
					return xerrors.Errorf("failed to read cid field t.Data: %w", err)
				}

				t.Data = c

			}
			// t.Prev (cid.Cid) (struct)
		case "prev":

			{

				b, err := cr.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := cr.UnreadByte(); err != nil {
						return err
					}

					c, err := cbg.ReadCid(cr)
					if err != nil {
						return xerrors.Errorf("failed to read cid field t.Prev: %w", err)
					}

					t.Prev = &c
				}

			}
			// t.Version (int64) (int64)
		case "version":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.Version = int64(extraI)
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}