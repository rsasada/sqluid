package backend

import (
	"os"
	"errors"
)

type Pager struct {
	file		*os.File
	FileLength	uint
	Pages		[TableMaxSize][]byte
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
	pager.FileLength = uint(fileSize)

	t.Pager = &pager
	return nil
}

func (t *Table) PagerFlush(pageNum uint, dataSize uint) error {

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

	_, err = t.Pager.file.Write(t.Pager.Pages[pageNum][:dataSize])
	if err != nil {
		return err
	}

	return nil
}

func (t *Table) PagerClose() error {

	pages := t.Pager.Pages
	rowSize :=  t.RowSize()
	RowsPerPage := PageSize / rowSize
	numPages := t.NumRows / RowsPerPage

	for i := uint(0); i < numPages; i ++ {

		if pages[i] == nil {
			continue
		}

		err := t.PagerFlush(i, PageSize)
		if err != nil {
			return err
		}

		t.Pager.Pages[i] = nil
	}

	leftOverRows := t.NumRows % RowsPerPage
	if leftOverRows != 0 {

		if pages[numPages] != nil {

			err := t.PagerFlush(numPages, leftOverRows * t.RowSize())
			if err != nil {
				return err
			}

			t.Pager.Pages[numPages] = nil
		}
	}

	t.Pager.file.Close()
	return nil
}

func (t *Table) SetPage(pageNum uint) ([]byte, error) {

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

		t.Pager.Pages[pageNum] = &newPage
	}

	return t.Pager.Pages[pageNum], nil
}
