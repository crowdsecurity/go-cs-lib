package yamlpatch

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Patcher struct {
	BaseFilePath  string
	PatchFilePath string
	quiet         bool
}

// Deprecated: use csyaml.NewPatcher instead.
func NewPatcher(filePath string, suffix string) *Patcher {
	return &Patcher{
		BaseFilePath:  filePath,
		PatchFilePath: filePath + suffix,
		quiet:         false,
	}
}

// SetQuiet sets the quiet flag, which will log as DEBUG_LEVEL instead of INFO.
func (p *Patcher) SetQuiet(quiet bool) {
	p.quiet = quiet
}

// read a single YAML file, check for errors (the merge package doesn't) then return the content as bytes.
func readYAML(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("while reading yaml file: %w", err)
	}

	var yamlMap map[any]any

	if err = yaml.Unmarshal(content, &yamlMap); err != nil {
		return nil, fmt.Errorf("%s: %w", filePath, err)
	}

	return content, nil
}

// MergedPatchContent reads a YAML file and, if it exists, its patch file,
// then merges them and returns it serialized.
func (p *Patcher) MergedPatchContent() ([]byte, error) {
	base, err := readYAML(p.BaseFilePath)
	if err != nil {
		return nil, err
	}

	over, err := readYAML(p.PatchFilePath)
	if errors.Is(err, os.ErrNotExist) {
		return base, nil
	}

	if err != nil {
		return nil, err
	}

	logf := log.Infof
	if p.quiet {
		logf = log.Debugf
	}

	logf("Loading yaml file: '%s' with additional values from '%s'", p.BaseFilePath, p.PatchFilePath)

	// strict mode true, will raise errors for duplicate map keys and
	// overriding with a different type
	patched, err := YAML([][]byte{base, over}, true)
	if err != nil {
		return nil, err
	}

	return patched.Bytes(), nil
}

// read multiple YAML documents inside a file, and writes them to a buffer
// separated by the appropriate '---' terminators.
func decodeDocuments(file *os.File, buf *bytes.Buffer, finalDashes bool) error {
	dec := yaml.NewDecoder(file)
	dec.SetStrict(true)

	dashTerminator := false

	for {
		yml := make(map[any]any)

		err := dec.Decode(&yml)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return fmt.Errorf("while decoding %s: %w", file.Name(), err)
		}

		docBytes, err := yaml.Marshal(&yml)
		if err != nil {
			return fmt.Errorf("while marshaling %s: %w", file.Name(), err)
		}

		if dashTerminator {
			buf.WriteString("---\n")
		}

		buf.Write(docBytes)

		dashTerminator = true
	}

	if dashTerminator && finalDashes {
		buf.WriteString("---\n")
	}

	return nil
}

// PrependedPatchContent collates the base .yaml file with the .yaml.patch, by putting
// the content of the patch BEFORE the base document. The result is a multi-document
// YAML in all cases, even if the base and patch files are single documents.
func (p *Patcher) PrependedPatchContent() ([]byte, error) {
	patchFile, err := os.Open(p.PatchFilePath)
	// optional file, ignore if it does not exist
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("while opening %s: %w", p.PatchFilePath, err)
	}

	var result bytes.Buffer

	if patchFile != nil {
		if err = decodeDocuments(patchFile, &result, true); err != nil {
			return nil, err
		}

		logf := log.Infof

		if p.quiet {
			logf = log.Debugf
		}

		logf("Prepending yaml: '%s' with '%s'", p.BaseFilePath, p.PatchFilePath)
	}

	baseFile, err := os.Open(p.BaseFilePath)
	if err != nil {
		return nil, fmt.Errorf("while opening %s: %w", p.BaseFilePath, err)
	}

	if err = decodeDocuments(baseFile, &result, false); err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}
