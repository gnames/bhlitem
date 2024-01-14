package gnbhl

import (
	"github.com/gnames/gnbhl/itembhl"
)

type GNBHL interface {
	Item(itemID uint) (itembhl.ItemBHL, error)
}
