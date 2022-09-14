package cachectl

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

func FsortPrintPagesStat(path string, re *regexp.Regexp) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	fis, _ := f.ReadDir(-1)
	err = f.Close()
	if err != nil {
		return err
	}

	finfoStore := []fileInfo{}
	for _, fi := range fis {
		fullpath := filepath.Join(path, fi.Name())
		if fi.IsDir() {
			fsortPrintPagesStat(fullpath, &finfoStore)
		} else {
			finfo, err := fi.Info()
			if err != nil {
				return err
			}
			finfoStore = append(finfoStore, fileInfo{finfo, fullpath})
		}
	}

	sort.Sort(ByFilesize(finfoStore))

    for _, fi := range finfoStore {
        if re.MatchString(fi.path) {
			PrintPagesStat(fi.path, fi.fi.Size())
		}
    }

	return nil
}

func fsortPrintPagesStat(searchpath string, store *[]fileInfo) error {
	f, err := os.Open(searchpath)
	if err != nil {
		return err
	}
	fis, _ := f.ReadDir(-1)
	err = f.Close()
	if err != nil {
		return err
	}

	for _, fi := range fis {
		fullpath := filepath.Join(searchpath, fi.Name())
		if fi.IsDir() {
			fsortPrintPagesStat(fullpath, store)
		} else {
			finfo, err := fi.Info()
			if err != nil {
				return err
			}
		    *store = append(*store, fileInfo{finfo, fullpath})
		}
	}

	return nil
}

type fileInfo struct {
	fi   os.FileInfo
	path string
}

type ByFilesize []fileInfo

func (fis ByFilesize) Len() int {
	return len(fis)
}

func (fis ByFilesize) Swap(i, j int) {
	fis[i], fis[j] = fis[j], fis[i]
}

func (fis ByFilesize) Less(i, j int) bool {
	return fis[i].fi.Size() < fis[j].fi.Size()
}
