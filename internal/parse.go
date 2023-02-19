package internal

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var ErrWalkError = errors.New("error walking directory")

type Parser interface {
	ParseSequences() (map[string]Sequence, error)
}

type FSParserOpts struct {
	Root string
}

func NewFSParser(opts FSParserOpts) *FSParser {
	return &FSParser{
		root: opts.Root,
	}
}

type FSParser struct {
	root string
}

func (f *FSParser) ParseSequences() (map[string]Sequence, error) {
	sequences := map[string]Sequence{}

	err := filepath.WalkDir(f.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), ".yaml") || strings.HasSuffix(d.Name(), ".yml") {
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading file: %w", err)
			}
			var seq Sequence
			err = yaml.Unmarshal(content, &seq)
			if err != nil {
				return fmt.Errorf("error unmarshalling: %w", err)
			}

			relPath, err := filepath.Rel(f.root, path)
			if err != nil {
				return fmt.Errorf("error computing relative path: %w", err)
			}
			sequences[relPath] = seq
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrWalkError, err)
	}
	return sequences, nil
}
