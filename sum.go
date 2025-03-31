package pkgx

import (
	"bytes"
	"maps"
	"os"
	"path/filepath"
	"slices"

	"golang.org/x/mod/sumdb/dirhash"
	gopkg "golang.org/x/tools/go/packages"
)

const SumFilename = "x.sum"

type Sum interface {
	Dir() string
	Bytes() []byte
	Save() error
	Hash(string) string
}

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

type sum struct {
	dir    string
	hashes map[string]string
}

func (s *sum) AddPackage(p *gopkg.Package) {
	if p.Dir != "" {
		h, _ := dirhash.HashDir(p.Dir, "", dirhash.Hash1)
		s.hashes[p.PkgPath] = h
		if p.Module != nil && s.dir == "" {
			s.dir = p.Module.Dir
		}
	}
}

func (s *sum) Dir() string {
	return s.dir
}

func (s *sum) Save() error {
	f, err := os.OpenFile(filepath.Join(s.dir, SumFilename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(s.Bytes())
	return err
}

func (s *sum) Bytes() []byte {
	b := bytes.NewBuffer(nil)

	for _, path := range slices.Sorted(maps.Keys(s.hashes)) {
		b.WriteString(path)
		b.WriteString(" ")
		b.WriteString(s.hashes[path])
		b.WriteString("\n")
	}

	return b.Bytes()
}

func (s *sum) Hash(path string) string {
	return s.hashes[path]
}
