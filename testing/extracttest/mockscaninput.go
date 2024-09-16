// Package extracttest provides structures to help create tabular tests for extractors.
package extracttest

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/osv-scalibr/extractor"
	"github.com/google/osv-scalibr/extractor/filesystem"
	scalibrfs "github.com/google/osv-scalibr/fs"
	"github.com/google/osv-scalibr/testing/fakefs"
)

// ScanInputMockConfig is used to quickly configure building a mock ScanInput
type ScanInputMockConfig struct {
	// Path of the file ScanInput will read, relative to the ScanRoot
	Path string
	// FakeScanRoot allows you to set a custom scanRoot, can be relative or absolute,
	// and will be translated to an absolute path
	FakeScanRoot string
	FakeFileInfo *fakefs.FakeFileInfo
}

// TestTableEntry is a entry to pass to ExtractionTester
type TestTableEntry struct {
	Name          string
	InputConfig   ScanInputMockConfig
	WantInventory []*extractor.Inventory
	WantErr       error
}

// CloseTestScanInput takes a scan input generated by GenerateScanInputMock
// and closes the associated file handle
func CloseTestScanInput(t *testing.T, si filesystem.ScanInput) {
	t.Helper()
	// If the Reader is not an io.Closer, then there is an implementation error and this should panic.
	err := si.Reader.(io.Closer).Close()
	if err != nil {
		t.Errorf("Close(): %v", err)
	}
}

// GenerateScanInputMock will try to open the file locally, and fail if the file doesn't exist
func GenerateScanInputMock(t *testing.T, config ScanInputMockConfig) filesystem.ScanInput {
	t.Helper()

	var scanRoot string
	if filepath.IsAbs(config.FakeScanRoot) {
		scanRoot = config.FakeScanRoot
	} else {
		workingDir, err := os.Getwd()
		if err != nil {
			t.Fatalf("Can't get working directory because '%s'", workingDir)
		}
		scanRoot = filepath.Join(workingDir, config.FakeScanRoot)
	}

	f, err := os.Open(filepath.Join(scanRoot, config.Path))
	if err != nil {
		t.Fatalf("Can't open test fixture '%s' because '%s'", config.Path, err)
	}
	info, err := f.Stat()
	if err != nil {
		t.Fatalf("Can't stat test fixture '%s' because '%s'", config.Path, err)
	}

	return filesystem.ScanInput{
		FS:     os.DirFS(scanRoot).(scalibrfs.FS),
		Path:   config.Path,
		Root:   scanRoot,
		Reader: f,
		Info:   info,
	}
}
