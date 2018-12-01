// Copyright 2018 the Service Broker Project Authors.
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

package brokerpak

import (
	getter "github.com/hashicorp/go-getter"
)

// fetchArchive uses go-getter to download archives. By default go-getter
// decompresses archives, so this configuration prevents that.
func fetchArchive(src, dest string) error {
	return (&getter.Client{
		Src:           src,
		Dst:           dest,
		Mode:          getter.ClientModeFile,
		Getters:       getter.Getters,
		Decompressors: map[string]getter.Decompressor{},
	}).Get()

	return getter.Get(src, dest)
}
