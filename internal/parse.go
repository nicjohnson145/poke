package internal

import (
	"bytes"
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

type SequenceMap map[string]Sequence

type Parser interface {
	Parse(path string) (SequenceMap, error)
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

func (f *FSParser) ParseSequences(root string) (SequenceMap, error) {
	sequences := SequenceMap{}

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
	decoder := yaml.NewDecoder(bytes.NewBuffer(content))
	decoder.KnownFields(true)
	var seq Sequence
	err = decoder.Decode(&seq)
	if err != nil {
		return Sequence{}, fmt.Errorf("error unmarshalling: %w", err)
	}

	seq.importedCalls = make(map[string]map[string]Call)

	seq.path = filepath.Dir(path)

	return seq, nil
}

func (f *FSParser) Parse(path string) (SequenceMap, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("error checking path: %w", err)
	}

	var sequences SequenceMap
	if info.IsDir() {
		sequences, err = f.ParseSequences(path)
		if err != nil {
			return nil, fmt.Errorf("error parsing directory: %w", err)
		}
	} else {
		seq, err := f.ParseSingleSequence(path)
		if err != nil {
			return nil, fmt.Errorf("error parsing file: %w", err)
		}
		sequences = SequenceMap{path: seq}
	}
	return sequences, nil
}
