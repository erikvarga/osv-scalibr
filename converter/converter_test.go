// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package converter_test

import (
	"math/rand"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/CycloneDX/cyclonedx-go"
	"github.com/spdx/tools-golang/spdx/v2/common"
	"github.com/spdx/tools-golang/spdx/v2/v2_3"
	"github.com/google/uuid"
	"github.com/google/osv-scalibr/converter"
	"github.com/google/osv-scalibr/extractor"
	"github.com/google/osv-scalibr/extractor/filesystem/language/python/wheelegg"
	"github.com/google/osv-scalibr/extractor/filesystem/sbom/spdx"
	"github.com/google/osv-scalibr/extractor/standalone/windows/dismpatch"
	"github.com/google/osv-scalibr/purl"
	scalibr "github.com/google/osv-scalibr"
)

func TestToSPDX23(t *testing.T) {
	// Make UUIDs deterministic
	uuid.SetRand(rand.New(rand.NewSource(1)))
	pipEx := wheelegg.New(wheelegg.DefaultConfig())

	testCases := []struct {
		desc       string
		scanResult *scalibr.ScanResult
		config     converter.SPDXConfig
		want       *v2_3.Document
	}{
		{
			desc: "Package with no custom config",
			scanResult: &scalibr.ScanResult{
				Inventories: []*extractor.Inventory{&extractor.Inventory{
					Name: "software", Version: "1.2.3", Extractor: pipEx,
				}},
			},
			want: &v2_3.Document{
				SPDXVersion:       "SPDX-2.3",
				DataLicense:       "CC0-1.0",
				SPDXIdentifier:    "DOCUMENT",
				DocumentName:      "SCALIBR-generated SPDX",
				DocumentNamespace: "https://spdx.google/81855ad8-681d-4d86-91e9-1e00167939cb",
				CreationInfo: &v2_3.CreationInfo{
					Creators: []common.Creator{
						common.Creator{
							CreatorType: "Tool",
							Creator:     "SCALIBR",
						},
					},
				},
				Packages: []*v2_3.Package{
					&v2_3.Package{
						PackageName:               "main",
						PackageSPDXIdentifier:     "SPDXRef-Package-main-52fdfc07-2182-454f-963f-5f0f9a621d72",
						PackageVersion:            "0",
						PackageDownloadLocation:   converter.NoAssertion,
						IsFilesAnalyzedTagPresent: false,
					},
					&v2_3.Package{
						PackageName:           "software",
						PackageSPDXIdentifier: "SPDXRef-Package-software-9566c74d-1003-4c4d-bbbb-0407d1e2c649",
						PackageVersion:        "1.2.3",
						PackageSupplier: &common.Supplier{
							Supplier:     converter.NoAssertion,
							SupplierType: converter.NoAssertion,
						},
						PackageDownloadLocation:   converter.NoAssertion,
						IsFilesAnalyzedTagPresent: false,
						PackageExternalReferences: []*v2_3.PackageExternalReference{
							&v2_3.PackageExternalReference{
								Category: "PACKAGE-MANAGER",
								RefType:  "purl",
								Locator:  "pkg:pypi/software@1.2.3",
							},
						},
					},
				},
				Relationships: []*v2_3.Relationship{
					&v2_3.Relationship{
						RefA: common.DocElementID{
							ElementRefID: "SPDXRef-Package-main-52fdfc07-2182-454f-963f-5f0f9a621d72",
						},
						RefB: common.DocElementID{
							ElementRefID: "SPDXRef-Package-software-9566c74d-1003-4c4d-bbbb-0407d1e2c649",
						},
						Relationship: "CONTAINS",
					},
					&v2_3.Relationship{
						RefA: common.DocElementID{
							ElementRefID: "SPDXRef-Package-software-9566c74d-1003-4c4d-bbbb-0407d1e2c649",
						},
						RefB: common.DocElementID{
							SpecialID: converter.NoAssertion,
						},
						Relationship: "CONTAINS",
					},
				},
			},
		},
		{
			desc: "Package with custom config",
			scanResult: &scalibr.ScanResult{
				Inventories: []*extractor.Inventory{&extractor.Inventory{
					Name: "software", Version: "1.2.3", Extractor: pipEx,
				}},
			},
			config: converter.SPDXConfig{
				DocumentName:      "Custom name",
				DocumentNamespace: "Custom namespace",
				Creators: []common.Creator{
					common.Creator{
						CreatorType: "Person",
						Creator:     "Custom creator",
					},
				},
			},
			want: &v2_3.Document{
				SPDXVersion:       "SPDX-2.3",
				DataLicense:       "CC0-1.0",
				SPDXIdentifier:    "DOCUMENT",
				DocumentName:      "Custom name",
				DocumentNamespace: "Custom namespace",
				CreationInfo: &v2_3.CreationInfo{
					Creators: []common.Creator{
						common.Creator{
							CreatorType: "Tool",
							Creator:     "SCALIBR",
						},
						common.Creator{
							CreatorType: "Person",
							Creator:     "Custom creator",
						},
					},
				},
				Packages: []*v2_3.Package{
					&v2_3.Package{
						PackageName:               "main",
						PackageSPDXIdentifier:     "SPDXRef-Package-main-6694d2c4-22ac-4208-a007-2939487f6999",
						PackageVersion:            "0",
						PackageDownloadLocation:   converter.NoAssertion,
						IsFilesAnalyzedTagPresent: false,
					},
					&v2_3.Package{
						PackageName:           "software",
						PackageSPDXIdentifier: "SPDXRef-Package-software-eb9d18a4-4784-445d-87f3-c67cf22746e9",
						PackageVersion:        "1.2.3",
						PackageSupplier: &common.Supplier{
							Supplier:     converter.NoAssertion,
							SupplierType: converter.NoAssertion,
						},
						PackageDownloadLocation:   converter.NoAssertion,
						IsFilesAnalyzedTagPresent: false,
						PackageExternalReferences: []*v2_3.PackageExternalReference{
							&v2_3.PackageExternalReference{
								Category: "PACKAGE-MANAGER",
								RefType:  "purl",
								Locator:  "pkg:pypi/software@1.2.3",
							},
						},
					},
				},
				Relationships: []*v2_3.Relationship{
					&v2_3.Relationship{
						RefA: common.DocElementID{
							ElementRefID: "SPDXRef-Package-main-6694d2c4-22ac-4208-a007-2939487f6999",
						},
						RefB: common.DocElementID{
							ElementRefID: "SPDXRef-Package-software-eb9d18a4-4784-445d-87f3-c67cf22746e9",
						},
						Relationship: "CONTAINS",
					},
					&v2_3.Relationship{
						RefA: common.DocElementID{
							ElementRefID: "SPDXRef-Package-software-eb9d18a4-4784-445d-87f3-c67cf22746e9",
						},
						RefB: common.DocElementID{
							SpecialID: converter.NoAssertion,
						},
						Relationship: "CONTAINS",
					},
				},
			},
		},
		{
			desc: "Package with invalid PURLs skipped",
			scanResult: &scalibr.ScanResult{
				Inventories: []*extractor.Inventory{
					// PURL field missing
					&extractor.Inventory{Extractor: pipEx},
					// No name
					&extractor.Inventory{
						Version: "1.2.3", Extractor: pipEx,
					},
					// No version
					&extractor.Inventory{
						Name: "software", Extractor: pipEx,
					},
				},
			},
			want: &v2_3.Document{
				SPDXVersion:       "SPDX-2.3",
				DataLicense:       "CC0-1.0",
				SPDXIdentifier:    "DOCUMENT",
				DocumentName:      "SCALIBR-generated SPDX",
				DocumentNamespace: "https://spdx.google/5fb90bad-b37c-4821-b6d9-5526a41a9504",
				CreationInfo: &v2_3.CreationInfo{
					Creators: []common.Creator{
						common.Creator{
							CreatorType: "Tool",
							Creator:     "SCALIBR",
						},
					},
				},
				Packages: []*v2_3.Package{&v2_3.Package{
					PackageName:               "main",
					PackageSPDXIdentifier:     "SPDXRef-Package-main-95af5a25-3679-41ba-a2ff-6cd471c483f1",
					PackageVersion:            "0",
					PackageDownloadLocation:   converter.NoAssertion,
					IsFilesAnalyzedTagPresent: false,
				}},
				Relationships: []*v2_3.Relationship{},
			},
		},
		{
			desc: "Invalid chars in package name replaced",
			scanResult: &scalibr.ScanResult{
				Inventories: []*extractor.Inventory{&extractor.Inventory{
					Name: "softw@re&", Version: "1.2.3", Extractor: pipEx,
				}},
			},
			want: &v2_3.Document{
				SPDXVersion:       "SPDX-2.3",
				DataLicense:       "CC0-1.0",
				SPDXIdentifier:    "DOCUMENT",
				DocumentName:      "SCALIBR-generated SPDX",
				DocumentNamespace: "https://spdx.google/0f070244-8615-4bda-8831-3f6a8eb668d2",
				CreationInfo: &v2_3.CreationInfo{
					Creators: []common.Creator{
						common.Creator{
							CreatorType: "Tool",
							Creator:     "SCALIBR",
						},
					},
				},
				Packages: []*v2_3.Package{
					&v2_3.Package{
						PackageName:               "main",
						PackageSPDXIdentifier:     "SPDXRef-Package-main-680b4e7c-8b76-4a1b-9d49-d4955c848621",
						PackageVersion:            "0",
						PackageDownloadLocation:   converter.NoAssertion,
						IsFilesAnalyzedTagPresent: false,
					},
					&v2_3.Package{
						PackageName:           "softw@re&",
						PackageSPDXIdentifier: "SPDXRef-Package-softw-re--6325253f-ec73-4dd7-a9e2-8bf921119c16",
						PackageVersion:        "1.2.3",
						PackageSupplier: &common.Supplier{
							Supplier:     converter.NoAssertion,
							SupplierType: converter.NoAssertion,
						},
						PackageDownloadLocation:   converter.NoAssertion,
						IsFilesAnalyzedTagPresent: false,
						PackageExternalReferences: []*v2_3.PackageExternalReference{
							&v2_3.PackageExternalReference{
								Category: "PACKAGE-MANAGER",
								RefType:  "purl",
								Locator:  "pkg:pypi/softw%40re%26@1.2.3",
							},
						},
					},
				},
				Relationships: []*v2_3.Relationship{
					&v2_3.Relationship{
						RefA: common.DocElementID{
							ElementRefID: "SPDXRef-Package-main-680b4e7c-8b76-4a1b-9d49-d4955c848621",
						},
						RefB: common.DocElementID{
							ElementRefID: "SPDXRef-Package-softw-re--6325253f-ec73-4dd7-a9e2-8bf921119c16",
						},
						Relationship: "CONTAINS",
					},
					&v2_3.Relationship{
						RefA: common.DocElementID{
							ElementRefID: "SPDXRef-Package-softw-re--6325253f-ec73-4dd7-a9e2-8bf921119c16",
						},
						RefB: common.DocElementID{
							SpecialID: converter.NoAssertion,
						},
						Relationship: "CONTAINS",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := converter.ToSPDX23(tc.scanResult, tc.config)
			// Can't mock time.Now() so skip verifying the timestamp.
			tc.want.CreationInfo.Created = got.CreationInfo.Created

			if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(v2_3.Package{})); diff != "" {
				t.Errorf("converter.ToSPDX23(%v): unexpected diff (-want +got):\n%s", tc.scanResult, diff)
			}
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}

func TestToCDX(t *testing.T) {
	// Make UUIDs deterministic
	uuid.SetRand(rand.New(rand.NewSource(1)))
	pipEx := wheelegg.New(wheelegg.DefaultConfig())
	defaultBOM := cyclonedx.NewBOM()

	testCases := []struct {
		desc       string
		scanResult *scalibr.ScanResult
		config     converter.CDXConfig
		want       *cyclonedx.BOM
	}{
		{
			desc: "Package with custom config",
			scanResult: &scalibr.ScanResult{
				Inventories: []*extractor.Inventory{&extractor.Inventory{
					Name: "software", Version: "1.2.3", Extractor: pipEx,
				}},
			},
			config: converter.CDXConfig{
				ComponentName:    "sbom-1",
				ComponentVersion: "1.0.0",
				Authors:          []string{"author"},
			},
			want: &cyclonedx.BOM{
				Metadata: &cyclonedx.Metadata{
					Component: &cyclonedx.Component{
						Name:    "sbom-1",
						Version: "1.0.0",
						BOMRef:  "52fdfc07-2182-454f-963f-5f0f9a621d72",
					},
					Authors: ptr([]cyclonedx.OrganizationalContact{{Name: "author"}}),
					Tools: &cyclonedx.ToolsChoice{
						Components: &[]cyclonedx.Component{
							{
								Type: cyclonedx.ComponentTypeApplication,
								Name: "SCALIBR",
								ExternalReferences: ptr([]cyclonedx.ExternalReference{
									{URL: "https://github.com/google/osv-scalibr", Type: cyclonedx.ERTypeWebsite},
								}),
							},
						},
					},
				},
				Components: ptr([]cyclonedx.Component{
					{
						BOMRef:     "9566c74d-1003-4c4d-bbbb-0407d1e2c649",
						Type:       "library",
						Name:       "software",
						Version:    "1.2.3",
						PackageURL: "pkg:pypi/software@1.2.3",
					},
				}),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := converter.ToCDX(tc.scanResult, tc.config)
			// Can't mock time.Now() so skip verifying the timestamp.
			tc.want.Metadata.Timestamp = got.Metadata.Timestamp
			// Auto-populated fields
			tc.want.XMLNS = defaultBOM.XMLNS
			tc.want.JSONSchema = defaultBOM.JSONSchema
			tc.want.BOMFormat = defaultBOM.BOMFormat
			tc.want.SpecVersion = defaultBOM.SpecVersion
			tc.want.Version = defaultBOM.Version

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("converter.ToCDX(%v): unexpected diff (-want +got):\n%s", tc.scanResult, diff)
			}
		})
	}
}

func TestToPURL(t *testing.T) {
	pipEx := wheelegg.New(wheelegg.DefaultConfig())
	tests := []struct {
		desc      string
		inventory *extractor.Inventory
		want      *purl.PackageURL
		wantErr   bool
		onGoos    string
	}{
		{
			desc: "Valid inventory extractor",
			inventory: &extractor.Inventory{
				Name:      "software",
				Version:   "1.0.0",
				Locations: []string{"/file1"},
				Extractor: pipEx,
			},
			want: &purl.PackageURL{
				Type:    purl.TypePyPi,
				Name:    "software",
				Version: "1.0.0",
			},
		},
		{
			desc: "Windows-only returns error on Linux",
			inventory: &extractor.Inventory{
				Name:      "irrelevant",
				Extractor: dismpatch.Extractor{},
				Locations: []string{"irrelevant"},
				Version:   "irrelevant",
			},
			wantErr: true,
			onGoos:  "linux",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			if tc.onGoos != "" && tc.onGoos != runtime.GOOS {
				t.Skipf("Skipping test on %s", runtime.GOOS)
			}

			got, err := converter.ToPURL(tc.inventory)
			if err != nil && !tc.wantErr || err == nil && tc.wantErr {
				t.Fatalf("converter.ToPURL(%v): %v", tc.inventory, err)
			}

			if tc.wantErr == true {
				return
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("converter.ToPURL(%v) returned unexpected diff (-want +got):\n%s", tc.inventory, diff)
			}
		})
	}
}

func TestToCPEs(t *testing.T) {
	tests := []struct {
		desc      string
		inventory *extractor.Inventory
		want      []string
		wantErr   bool
	}{
		{
			desc: "Valid fileststem extractor",
			inventory: &extractor.Inventory{
				Name: "nginx",
				Metadata: &spdx.Metadata{
					CPEs: []string{"cpe:2.3:a:nginx:nginx:1.21.1"},
				},
				Extractor: &spdx.Extractor{},
			},
			want: []string{"cpe:2.3:a:nginx:nginx:1.21.1"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := converter.ToCPEs(tc.inventory)
			if err != nil && !tc.wantErr || err == nil && tc.wantErr {
				t.Fatalf("converter.ToCPEs(%v): %v", tc.inventory, err)
			}

			if tc.wantErr == true {
				return
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("converter.ToCPEs(%v) returned unexpected diff (-want +got):\n%s", tc.inventory, diff)
			}
		})
	}
}
