package itembhl

import (
	"github.com/gnames/gnbhl/ent/chunkbhl"
	"github.com/gnames/gnbhl/ent/pagebhl"
)

// ItemBHL is a BHL item.
type ItemBHL interface {
	// Pages returns a list of pages in the item.
	// Pages are sorted by their appearance in the item.
	Pages() []*pagebhl.PageBHL

	// Text returns the text of the item.
	Text() string

	// Page returns a page by its ID. It also returns the index of the page.
	// It returns an error if the page is not found.
	Page(id uint) (*pagebhl.PageBHL, int, error)

	// PageByOffset returns a page by a given offset.
	// The offset is the number of characters from the beginning of the item.
	// The offset should be located within the page.
	// It also returns the index of the page.
	// It returns an error if the page is not found.
	PageByOffset(offset uint) (*pagebhl.PageBHL, int, error)

	// PageText returns the text of a page by its ID.
	PageText(id uint) (string, error)

	// Chunk returns a chunk of text by its start and end offsets.
	Chunk(start, end uint) (*chunkbhl.ChunkBHL, error)
}
