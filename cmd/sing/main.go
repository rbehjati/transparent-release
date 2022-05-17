// Copyright 2022 The Project Oak Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Run with `go run cmd/sing/main.go -- provenance_path ./schema/amber-slsa-buildtype/v1/example.json`

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	// Enable the github OIDC auth provider.
	_ "github.com/sigstore/cosign/pkg/providers/filesystem"
	_ "github.com/sigstore/cosign/pkg/providers/github"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
	"github.com/slsa-framework/slsa-github-generator/signing/sigstore"
)

func check(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, e.Error())
		os.Exit(1)
	}
}

func signProvenance(path string) error {
	// Generate the provenance.
	ctx := context.Background()

	provenanceBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read the provenance file: %v", err)
	}
	var provenance intoto.Statement
	err = json.Unmarshal(provenanceBytes, &provenance)

	// Sign the provenance.
	s := sigstore.NewDefaultFulcio()
	att, err := s.Sign(ctx, &provenance)
	if err != nil {
		return err
	}

	attBytes := att.Bytes()
	err = ioutil.WriteFile("intoto.jsonl", attBytes, 0600)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	provenancePath := flag.String("provenance_path", "provenance.json",
		"Path to the input provenance file.")

	flag.Parse()

	absProvenancePath, err := filepath.Abs(*provenancePath)

	err = signProvenance(absProvenancePath)
	check(err)
}
