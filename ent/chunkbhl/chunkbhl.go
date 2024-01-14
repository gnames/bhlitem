package chunkbhl

import "github.com/gnames/gnbhl/ent/pagebhl"

// ChunkBHL is a chunk of text from a BHL item.
type ChunkBHL struct {
	// Start is the offset of the first character of the chunk.
	Start uint

	// End is the offset of the last character of the chunk.
	End uint

	// PageIdxStart is the index of the first page in the chunk.
	PageIdxStart int

	// PageIdxEnd is the index of the last page in the chunk.
	PageIdxEnd int

	// Pages is a list of pages in the chunk.
	Pages []*pagebhl.PageBHL

	// Text is the text of the chunk.
	Text string
}
