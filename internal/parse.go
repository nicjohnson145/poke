package internal

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

var ErrWalkError = errors.New("error walking directory")

type Parser interface {
	ParseSequences(root string) (map[string]Sequence, error)
	ParseSingleSequence(path string) (Sequence, error)
}

type FSParserOpts struct {
	Logger zerolog.Logger
}

func NewFSParser(opts FSParserOpts) *FSParser {
	return &FSParser{
		log: opts.Logger,
	}
}

var _ Parser = (*FSParser)(nil)

type FSParser struct {
	log zerolog.Logger
}

func (f *FSParser) ParseSequences(root string) (map[string]Sequence, error) {
	sequences := map[string]Sequence{}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), ".yaml") || strings.HasSuffix(d.Name(), ".yml") {
			seq, err := f.ParseSingleSequence(path)
			if err != nil {
				return fmt.Errorf("error parsing sequence: %w", err)
			}
			relPath, err := filepath.Rel(root, path)
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

func (f *FSParser) ParseSingleSequence(path string) (Sequence, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Sequence{}, fmt.Errorf("error reading file: %w", err)
	}
	var seq Sequence
	err = yaml.Unmarshal(content, &seq)
	if err != nil {
		return Sequence{}, fmt.Errorf("error unmarshalling: %w", err)
	}

	return seq, nil
}
