package fsio

import (
	"cmp"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/gnames/gnbhl/ent/pagebhl"
	"github.com/gnames/gnsys"
)

func GetPath(rootPath string, itemID uint) (string, error) {
	sID := fmt.Sprintf("%06d", itemID)
	if len(sID) != 6 {
		return "", fmt.Errorf("ItemID %d is not correct", itemID)
	}
	path := filepath.Join(rootPath, sID[0:3], sID)
	exists, empty, err := gnsys.DirExists(path)
	if err != nil || !exists || empty {
		err = fmt.Errorf("path '%s' does not exist or is empty: %w", path, err)
		return "", err
	}
	return path, nil
}

func GetPages(
	path string, itemID uint,
) (byID, bySeq []*pagebhl.PageBHL, err error) {
	var es []os.DirEntry
	es, err = os.ReadDir(path)
	if err != nil {
		return nil, nil, err
	}
	for i := range es {
		name := es[i].Name()
		if filepath.Ext(name) != ".txt" {
			continue
		}
		pageID, seqNum, err := nameToPageID(es[i].Name())
		if err != nil {
			return nil, nil, err
		}
		filePath := filepath.Join(path, name)
		bs, err := os.ReadFile(filePath)
		if err != nil {
			return nil, nil, err
		}
		txt := string(bs)
		page := pagebhl.PageBHL{
			ID:       uint(pageID),
			ItemID:   itemID,
			SeqNum:   uint(seqNum),
			FileName: name,
			Text:     txt,
			Length:   uint(len(txt)),
		}
		byID = append(byID, &page)
		bySeq = append(bySeq, &page)
	}
	slices.SortFunc(byID, func(a, b *pagebhl.PageBHL) int {
		return cmp.Compare(a.ID, b.ID)
	})
	slices.SortFunc(bySeq, func(a, b *pagebhl.PageBHL) int {
		return cmp.Compare(a.SeqNum, b.SeqNum)
	})

	offset := uint(0)
	for i := range bySeq {
		bySeq[i].Offset = offset
		offset += bySeq[i].Length
	}

	return byID, bySeq, nil
}

func nameToPageID(fName string) (id, seqNum int, err error) {
	fName = strings.TrimSuffix(fName, ".txt")
	es := strings.Split(fName, "-")
	if len(es) != 3 {
		return 0, 0, fmt.Errorf("bad file name %s", fName)
	}
	id, err = strconv.Atoi(es[1])
	if err != nil {
		return 0, 0, err
	}
	seqNum, err = strconv.Atoi(es[2])
	if err != nil {
		return 0, 0, err
	}

	return id, seqNum, nil
}
