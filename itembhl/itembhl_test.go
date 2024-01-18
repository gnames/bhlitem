package itembhl_test

import (
	"testing"

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
	itemID := uint(100100)
	res, err = itembhl.New("bad_path", itemID)
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
	start := 13000
	end := 14000
	cnk, err := res.Chunk(start, end)
	assert.Nil(err)
	assert.Equal(uint(start), cnk.Start)
	assert.Equal(uint(end), cnk.End)
	assert.Equal(1000, len(cnk.Text))
	assert.Equal(2, len(cnk.Pages))
	assert.Equal(uint(17), cnk.Pages[0].SeqNum)
	assert.Equal(uint(18), cnk.Pages[1].SeqNum)

	cnk, err = res.Chunk(start, 100_000_000)
	assert.Nil(err)
	assert.Equal(13000, int(cnk.Start))
	assert.Equal(50998, int(cnk.End))

	cnk, err = res.Chunk(-50, 100_000_000)
	assert.Nil(err)
	assert.Equal(0, int(cnk.Start))
	assert.Equal(50998, int(cnk.End))
}

func getItem(path string, t *testing.T) itembhl.ItemBHL {
	assert := assert.New(t)
	if itm != nil {
		return itm
	}
	itemID := uint(100100)
	res, err := itembhl.New(path, itemID)
	assert.Nil(err)
	return res
}
