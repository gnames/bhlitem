package itembhl_test

import (
	"testing"

	"github.com/gnames/gnbhl/config"
	"github.com/gnames/gnbhl/itembhl"
	"github.com/stretchr/testify/assert"
)

var (
	itm  itembhl.ItemBHL
	path = "../testdata/bhl"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	res := getItem(path, t)
	assert.Implements((*itembhl.ItemBHL)(nil), res)

	var err error
	cfg := config.New(config.OptPath("bad_path"))
	itemID := uint(100100)
	res, err = itembhl.New(cfg, itemID)
	assert.NotNil(err)
	assert.Nil(res)
}

func TestPages(t *testing.T) {
	assert := assert.New(t)
	res := getItem(path, t)
	pages := res.Pages()
	assert.Equal(68, len(pages))
	assert.Equal(pages[0].ID, uint(32076427))
	assert.Equal(pages[len(pages)-1].ID, uint(32076389))
}

func TestText(t *testing.T) {
	assert := assert.New(t)
	res := getItem(path, t)
	txt := res.Text()
	assert.Greater(len(txt), 1000)
}

func TestPage(t *testing.T) {
	assert := assert.New(t)
	res := getItem(path, t)
	var pgNum uint = 32076400
	pg, _, err := res.Page(pgNum)
	assert.Nil(err)
	assert.Equal(pgNum, pg.ID)
	assert.Equal("100100-32076400-0028.txt", pg.FileName)
	assert.Equal(uint(28), pg.SeqNum)

	pg, _, err = res.Page(uint(42))
	assert.NotNil(err)
	assert.Nil(pg)
}

func TestOffset(t *testing.T) {
	assert := assert.New(t)
	res := getItem(path, t)
	var offset uint = 3000
	pg, _, err := res.PageByOffset(offset)
	assert.Equal("100100-32076421-0007.txt", pg.FileName)
	assert.Nil(err)
	assert.Equal(uint(7), pg.SeqNum)
}

func TestChunk(t *testing.T) {
	assert := assert.New(t)
	res := getItem(path, t)
	start := uint(13000)
	end := uint(14000)
	chk, err := res.Chunk(start, end)
	assert.Nil(err)
	assert.Equal(start, chk.Start)
	assert.Equal(end, chk.End)
	assert.Equal(1000, len(chk.Text))
	assert.Equal(2, len(chk.Pages))
	assert.Equal(uint(17), chk.Pages[0].SeqNum)
	assert.Equal(uint(18), chk.Pages[1].SeqNum)

	chk, err = res.Chunk(end, start)
	assert.NotNil(err)
	assert.Nil(chk)
	_, err = res.Chunk(start, 100_000_000)
	assert.NotNil(err)
}

func getItem(path string, t *testing.T) itembhl.ItemBHL {
	assert := assert.New(t)
	if itm != nil {
		return itm
	}
	cfg := config.New(config.OptPath(path))
	itemID := uint(100100)
	res, err := itembhl.New(cfg, itemID)
	assert.Nil(err)
	return res
}
