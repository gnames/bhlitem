package itembhl

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/gnames/gnbhl/ent/chunkbhl"
	"github.com/gnames/gnbhl/ent/pagebhl"
	"github.com/gnames/gnbhl/internal/fsio"
	"github.com/gnames/gnlib"
)

type itembhl struct {
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

// New creates an instance of ItemBHL.
func New(rootPath string, itemID uint) (ItemBHL, error) {
	path, err := fsio.GetPath(rootPath, itemID)
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
		id:         itemID,
		path:       path,
		pagesByID:  byID,
		pagesBySeq: bySeq,
		text:       text,
		length:     uint(len(text)),
	}

	return res, nil
}

// Pages returns a list of pages in the item. The pages are sorted by their
// appearance in the item.
func (itm *itembhl) Pages() []*pagebhl.PageBHL {
	return itm.pagesBySeq
}

// Text returns the text of the whole item.
func (itm *itembhl) Text() string {
	return itm.text
}

// Page returns a page by its ID. It also returns the index of the page in the
// pages' sequence.
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

// PageByOffset returns a page that contains a given offset. The offset is the
// number of the UTF-8 characters from the beginning of the item. The offset
// should be located within the page. It also returns the index of the page in
// the pages' sequence.
func (itm *itembhl) PageByOffset(offset uint) (*pagebhl.PageBHL, int, error) {
	idx, found := slices.BinarySearchFunc(
		itm.pagesBySeq, offset,
		func(a *pagebhl.PageBHL, b uint) int {
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

// PageText returns the text of a page by its ID.
func (itm *itembhl) PageText(id uint) (string, error) {
	pg, _, err := itm.Page(id)
	if err != nil {
		return "", err
	}
	return pg.Text, nil
}

// Chunk returns a chunk of text by its start and end offsets. The offset is
// the number of UTF-8 characters from the beginning of the item.
// If the chunk is out of bounds, the bounds are adjusted to the beginning or
// the end of the item.
func (itm *itembhl) Chunk(start, end int) (*chunkbhl.ChunkBHL, error) {
	if start > end {
		return nil, fmt.Errorf("start offset %d is greater than end offset %d", start, end)
	}
	l := int(itm.length)
	if end > l {
		end = l
	}

	if start < 0 {
		start = 0
	}

	_, idxStart, err := itm.PageByOffset(uint(start))
	if err != nil {
		return nil, err
	}
	_, idxEnd, err := itm.PageByOffset(uint(end))
	if err != nil {
		return nil, err
	}
	chunk := &chunkbhl.ChunkBHL{
		Start:        uint(start),
		End:          uint(end),
		PageIdxStart: idxStart,
		PageIdxEnd:   idxEnd,
		Pages:        itm.pagesBySeq[idxStart : idxEnd+1],
		Text:         itm.Text()[start:end],
	}
	return chunk, nil
}
