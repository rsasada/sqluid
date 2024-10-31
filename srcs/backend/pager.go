package backend

import (
	"os"
	"fmt"
	"errors"
	"encoding/binary"
)

type Pager struct {
	file		*os.File
	FileLength	uint32
	Pages		[TableMaxSize][]byte
	NumPages	uint32
}

func (t *Table)PagerOpen(tableName string) error {
	
	pager := Pager{}
	filepath := tableName + ".idb"
	
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	pager.file = file
	
	fileInfo, err := file.Stat()
    if err != nil {
		return err
    }

	fileSize := fileInfo.Size()
	if fileSize < 0 {
		return errors.New("failed to get the size of the '.idp' file")
	}
	pager.FileLength = uint32(fileSize)

	pager.NumPages = pager.FileLength / uint32(PageSize)
	if pager.FileLength % PageSize != 0 {
		fmt.Println("Jeeez!! your DB file get fucked up,,,")
	}

	t.RootPageNum = 0

	if pager.NumPages == 0 {
		pager.Pages[0], err = t.SetPage(0)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint32(pager.Pages[0][numCellsOffset:numCellsOffset+4], 0)
	}

	t.Pager = &pager
	return nil
}

func (t *Table) PagerFlush(pageNum uint32) error {

	if pageNum > TableMaxSize {
		return errors.New("Tried to flush page number out of bounds")
	}

	if t.Pager.Pages[pageNum] == nil {
		return errors.New("Tried to flush nil page")
	}

	_, err := t.Pager.file.Seek(int64(PageSize * pageNum), 0)
	if err != nil {
		return err
	}

	_, err = t.Pager.file.Write(t.Pager.Pages[pageNum])
	if err != nil {
		return err
	}

	return nil
}

func (t *Table) PagerClose() error {

	pages := t.Pager.Pages

	for i := uint32(0); i < t.Pager.NumPages; i ++ {

		if pages[i] == nil {
			continue
		}

		err := t.PagerFlush(i)
		if err != nil {
			return err
		}

		t.Pager.Pages[i] = nil
	}

	t.Pager.file.Close()
	return nil
}

func (t *Table) SetPage(pageNum uint32) ([]byte, error) {

	if pageNum > TableMaxSize {
		return nil, errors.New("Tried to fetch page number out of bounds")
	}
	
	if t.Pager.Pages[pageNum] == nil {

		newPage := make([]byte, PageSize)
		numPages := t.Pager.FileLength / PageSize

		if t.Pager.FileLength % PageSize != 0 {
			numPages += 1
		}

		if pageNum <= numPages {
			_, err := t.Pager.file.Seek(int64(PageSize * pageNum), 0)
			if err != nil {
				return nil, err
			}

			_, err = t.Pager.file.Read(newPage)
			if err != nil {
				return nil, err
			}
		}

		t.Pager.Pages[pageNum] = newPage

		if pageNum >= t.Pager.NumPages {
			t.Pager.NumPages = pageNum + 1
		}
	}

	return t.Pager.Pages[pageNum], nil
}
