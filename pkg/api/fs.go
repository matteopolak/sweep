package api

import (
	"io/fs"
	"os"
	"path"
)

// Struct that represents one predicate-folder pair, which allows
// us to only search for certain folders within a given directory
// when a certain file exists (e.g. package.json with node_modules)
type Predicate struct {
	Predicate func(string) bool
	Folder    string
	Match     bool
}

// Deletes an entire directory
func DeleteDir(dir string) error {
	return os.RemoveAll(dir)
}

// The state of a directory iterator, which keeps track of the current
// directory content, and the next directories that can be searched
type DirIterator struct {
	CurrentDirContent []fs.DirEntry
	CurrentDir        string
	NextDirs          []string
	Predicates        []Predicate
	PredicateIndex    int
	Ignore            []string
}

func (d *DirIterator) Next() (string, error) {
	// If the current directory content is empty, fetch the next directory contents that aren't empty
	if len(d.CurrentDirContent) == 0 {
		i := 0

		// iterate over NextDirs, and set CurrentDirContent to the first directory
		// that has content
		// then remove all directories from 0 to i
		for i < len(d.NextDirs) {
			// get the directory content
			dirContent, err := os.ReadDir(d.NextDirs[i])

			// if there's an error, return it
			if err != nil {
				return "", err
			}

			// if the directory content is not empty, set CurrentDirContent to it
			// and break out of the loop
			if len(dirContent) > 0 {
				// filter the directory content and check again  if it's empty
				filteredDirContent := []fs.DirEntry{}

				for _, dir := range dirContent {
					// if the directory is in the ignore list, skip it
					ignore := false

					for _, ignored := range d.Ignore {
						if ignored == dir.Name() {
							ignore = true
							break
						}
					}

					if !ignore {
						filteredDirContent = append(filteredDirContent, dir)
					}
				}

				if len(filteredDirContent) > 0 {
					d.CurrentDirContent = dirContent
					break
				}
			}

			i++
		}

		// if i is equal to the length of NextDirs, then we didn't find any directories
		// with content, so we're done
		if i == len(d.NextDirs) {
			return "", nil
		}

		d.CurrentDir = d.NextDirs[i]
		d.NextDirs = d.NextDirs[i+1:]

		// set Match to false
		for i := range d.Predicates {
			d.Predicates[i].Match = false
		}

		// Iterate over the predicates, and set MatchFolder to true if the Folder
		// exists in the current directory content
		for i := range d.Predicates {
			for _, dir := range d.CurrentDirContent {
				if !d.Predicates[i].Match && d.Predicates[i].Predicate(dir.Name()) {
					d.Predicates[i].Match = true
				}
			}
		}

		// iterate over the directory content, if the entry is a directory and didn't match a predicate,
		// add it to the NextDirs
		for _, dir := range d.CurrentDirContent {
			if dir.IsDir() {
				// iterate over the predicates, and if the directory name doesn't match any of them,
				// add it to NextDirs
				matched := false

				for _, predicate := range d.Predicates {
					if predicate.Folder == dir.Name() {
						matched = true
						break
					}
				}

				if !matched {
					// append full path
					d.NextDirs = append(d.NextDirs, path.Join(d.CurrentDir, dir.Name()))
				}
			}
		}

		// Then, reset PredicateIndex to 0
		d.PredicateIndex = 0
	}

	// Now, iterate over the predicates again, and if Match is true, check the predicate
	// and yield the folder if it's true. Then, set MatchFolder back to false
	// Note that we should start at PredicateIndex
	for i := d.PredicateIndex; i < len(d.Predicates); i++ {
		if d.Predicates[i].Match {
			for _, dir := range d.CurrentDirContent {
				if d.Predicates[i].Folder == dir.Name() {
					d.PredicateIndex = i + 1
					d.Predicates[i].Match = false

					return path.Join(d.CurrentDir, dir.Name()), nil
				}
			}
		}
	}

	// If we get here, we didn't find any folders that matched the predicates
	// so just empty out the current directory content, and call Next
	d.CurrentDirContent = []fs.DirEntry{}

	return d.Next()
}

// Travserses a directory, and yields all files that should be deleted,
// based on the predicate-folder pairs passed in
func NewDirIterator(dir string, pairs []Predicate, ignore []string) *DirIterator {
	// Create the iterator
	iterator := &DirIterator{
		CurrentDirContent: []fs.DirEntry{},
		NextDirs:          []string{dir},
		Predicates:        pairs,
		Ignore:            ignore,
		PredicateIndex:    0,
		CurrentDir:        "",
	}

	return iterator
}

func GetDirSize(dir string) int64 {
	size := int64(0)

	// get the directory content
	dirContent, err := os.ReadDir(dir)

	// if there's an error, return it
	if err != nil {
		return 0
	}

	// iterate over the directory content, and if the entry is a file,
	// add its size to the total size
	for _, dir := range dirContent {
		if !dir.IsDir() {
			fileInfo, err := dir.Info()

			if err != nil {
				continue
			}

			size += fileInfo.Size()
		}
	}

	return size
}
