package pagebhl

// PageBHL represents a page in an BHL item.
type PageBHL struct {
	// ID is the unique BHL-assigned ID of the page.
	ID uint `json:"id"`

	// ItemID is the ID of the item that contains the page.
	ItemID uint `json:"itemId"`

	// FileName is the name of the file that contains the page.
	FileName string `json:"fileName"`

	// Number is the page number of the page.
	Number uint `json:"number,omitempty"`

	// SecNum is the page number given to a page by Arxive.
	SeqNum uint `json:"seqNum"`

	// Offset is the offset of the page in the item (in runes).
	Offset uint `json:"offset"`

	// Length is the length of the page (in runes).
	Length uint `json:"length"`

	// Text is the text of the page.
	Text string `json:"text"`
}
