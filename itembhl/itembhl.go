package itembhl

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/gnames/gnbhl/config"
	"github.com/gnames/gnbhl/ent/chunkbhl"
	"github.com/gnames/gnbhl/ent/pagebhl"
	"github.com/gnames/gnbhl/internal/fsio"
	"github.com/gnames/gnlib"
)

type itembhl struct {
	// cfg is the configuration of gnbhl
	cfg config.Config

	// id is the unique BHL-assigned ID of the item.
	id uint

	// path is the filesystem path to the item directory
	path string

	// text is the content of the item.
	text string

	// length is the length of the item (in runes).
	length uint

	// pagesByID is a list of pagesByID in the item.
	// The pagesByID should be sorted by ID.
	pagesByID []*pagebhl.PageBHL

	pagesBySeq []*pagebhl.PageBHL
}

func New(cfg config.Config, itemID uint) (ItemBHL, error) {
	path, err := fsio.GetPath(cfg.Path, itemID)
	if err != nil {
		return nil, err
	}
	byID, bySeq, err := fsio.GetPages(path, itemID)
	if err != nil {
		return nil, err
	}

	txts := gnlib.Map(bySeq, func(p *pagebhl.PageBHL) string {
		return p.Text
	})
	text := strings.Join(txts, "")

	res := &itembhl{
		cfg:        cfg,
		id:         itemID,
		path:       path,
		pagesByID:  byID,
		pagesBySeq: bySeq,
		text:       text,
		length:     uint(len(text)),
	}

	return res, nil
}

func (itm *itembhl) Pages() []*pagebhl.PageBHL {
	return itm.pagesBySeq
}

func (itm *itembhl) Text() string {
	return itm.text
}

func (itm *itembhl) Page(id uint) (*pagebhl.PageBHL, int, error) {
	idx, found := slices.BinarySearchFunc(
		itm.pagesByID, id,
		func(a *pagebhl.PageBHL, b uint) int {
			return cmp.Compare(a.ID, b)
		},
	)

	if !found {
		return nil, -1, fmt.Errorf("page with ID %d not found", id)
	}
	return itm.pagesByID[idx], idx, nil
}

func (itm *itembhl) PageByOffset(offset uint) (*pagebhl.PageBHL, int, error) {
	idx, found := slices.BinarySearchFunc(
		itm.pagesBySeq, offset,
		func(a *pagebhl.PageBHL, b uint) int {
			fmt.Printf("SEC: %d, a.Offset: %d, b: %d\n", a.SeqNum, a.Offset, b)
			if a.Offset > b {
				return 1
			}
			if a.Offset <= b && a.Offset+a.Length > b {
				return 0
			}
			return -1
		},
	)
	if !found {
		return nil, -1, fmt.Errorf("page with offset %d not found", offset)
	}
	return itm.pagesBySeq[idx], idx, nil
}

func (itm *itembhl) PageText(id uint) (string, error) {
	pg, _, err := itm.Page(id)
	if err != nil {
		return "", err
	}
	return pg.Text, nil
}

func (itm *itembhl) Chunk(start, end uint) (*chunkbhl.ChunkBHL, error) {
	if start > end {
		return nil, fmt.Errorf("start offset %d is greater than end offset %d", start, end)
	}
	if end > itm.length {
		return nil, fmt.Errorf("end offset %d is greater than item length %d", end, itm.length)
	}
	_, idxStart, err := itm.PageByOffset(start)
	if err != nil {
		return nil, err
	}
	_, idxEnd, err := itm.PageByOffset(end)
	if err != nil {
		return nil, err
	}
	chunk := &chunkbhl.ChunkBHL{
		Start:        start,
		End:          end,
		PageIdxStart: idxStart,
		PageIdxEnd:   idxEnd,
		Pages:        itm.pagesBySeq[idxStart : idxEnd+1],
		Text:         itm.text[start:end],
	}
	return chunk, nil
}
