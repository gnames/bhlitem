package gnbhl

import (
	"github.com/gnames/gnbhl/config"
	"github.com/gnames/gnbhl/itembhl"
)

type gnbhl struct {
	cfg config.Config
}

func New(cfg config.Config) GNBHL {
	res := gnbhl{
		cfg: cfg,
	}
	return &res
}

func (g *gnbhl) Item(itemID uint) (itembhl.ItemBHL, error) {
	return itembhl.New(g.cfg.Path, itemID)
}
