package internal

import (
	"bytes"
	"maps"
	"os"
	"path/filepath"
	"slices"

	"golang.org/x/mod/sumdb/dirhash"
	gopkg "golang.org/x/tools/go/packages"
)

const SumFilename = "go.xsum"

func LoadSumFile(m *gopkg.Module) Sum {
	if m == nil || m.Dir == "" {
		return nil
	}
	s := &sum{
		dir:    m.Dir,
		hashes: make(map[string]string),
	}

	data, err := os.ReadFile(filepath.Join(m.Dir, SumFilename))
	if err != nil {
		return nil
	}

	for line := range bytes.Lines(data) {
		parts := bytes.Fields(line)
		if len(parts) == 2 {
			s.hashes[string(parts[0])] = string(parts[1])
		}
	}

	return s
}

// Sum helps to calculate module's sum file
type Sum interface {
	// Dir returns module's source dir
	Dir() string
	// Add adds hash of package
	Add(*gopkg.Package)
	// Save saves xsum file to module's source dir
	Save() error
	// Hash returns hash of package by package path
	Hash(string) string
}

func NewSum(dir string) Sum {
	return &sum{dir: dir, hashes: make(map[string]string)}
}

type sum struct {
	// dir module source dir
	dir string
	// hashes of packages
	hashes map[string]string
}

func (s *sum) Dir() string { return s.dir }

func (s *sum) Add(p *gopkg.Package) {
	if _, ok := s.hashes[p.ID]; !ok {
		h, _ := dirhash.HashDir(p.Dir, "", dirhash.Hash1)
		s.hashes[p.ID] = h
	}
}

func (s *sum) Hash(path string) string { return s.hashes[path] }

func (s *sum) Save() error {
	b := bytes.NewBuffer(nil)

	for _, path := range slices.Sorted(maps.Keys(s.hashes)) {
		b.WriteString(path)
		b.WriteString(" ")
		b.WriteString(s.hashes[path])
		b.WriteString("\n")
	}

	f, err := os.OpenFile(
		filepath.Join(s.dir, SumFilename),
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		os.ModePerm,
	)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	_, err = f.Write(b.Bytes())
	return err
}
