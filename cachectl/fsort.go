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
	fis, _ := f.Readdir(-1)
	err = f.Close()
	if err != nil {
		return err
	}
	sort.Sort(ByFilesize(fis))

	for _, fi := range fis {
		fullpath, err := filepath.Abs(fi.Name())
		if err != nil {
			return err
		}
		if re.MatchString(fullpath) {
			PrintPagesStat(fullpath, fi.Size())
		}
	}

	return nil
}

type ByFilesize []os.FileInfo

func (fis ByFilesize) Len() int {
	return len(fis)
}

func (fis ByFilesize) Swap(i, j int) {
	fis[i], fis[j] = fis[j], fis[i]
}

func (fis ByFilesize) Less(i, j int) bool {
	return fis[i].Size() < fis[j].Size()
}
