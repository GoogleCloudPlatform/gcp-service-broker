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
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils/stream"
	getter "github.com/hashicorp/go-getter"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// fetchArchive uses go-getter to download archives. By default go-getter
// decompresses archives, so this configuration prevents that.
func fetchArchive(src, dest string) error {
	return newFileGetterClient(src, dest).Get()
}

// fetchBrokerpak downloads a local or remote brokerpak; brokerpaks can be
// fetched remotely using the gs:// prefix which will load them from a
// Cloud Storage bucket with the broker's credentials.
// Relative paths are resolved relative to the executable.
func fetchBrokerpak(src, dest string) error {
	execWd := filepath.Dir(os.Args[0])
	execDir, err := filepath.Abs(execWd)
	if err != nil {
		return fmt.Errorf("couldn't turn dir %q into abs path: %v", execWd, err)
	}

	client := newFileGetterClient(src, dest)
	client.Getters["gs"] = &gsGetter{}
	client.Pwd = execDir

	return client.Get()
}

func defaultGetters() map[string]getter.Getter {
	getters := map[string]getter.Getter{}
	for k, g := range getter.Getters {
		getters[k] = g
	}

	return getters
}

// newFileGetterClient creates a new client that will fetch a single file,
// with the default set of getters and will NOT automatically decompress it.
func newFileGetterClient(src, dest string) *getter.Client {
	return &getter.Client{
		Src: src,
		Dst: dest,

		Mode:          getter.ClientModeFile,
		Getters:       defaultGetters(),
		Decompressors: map[string]getter.Decompressor{},
	}
}

// gsGetter is a go-getter that works on Cloud Storage using the broker's
// service account. It's incomplete in that it doesn't support directories.
type gsGetter struct{}

// ClientMode is unsupported for gsGetter.
func (g *gsGetter) ClientMode(u *url.URL) (getter.ClientMode, error) {
	return getter.ClientModeInvalid, errors.New("mode is not supported for this client")
}

// Get clones a remote destination to a local directory.
func (g *gsGetter) Get(dst string, u *url.URL) error {
	return errors.New("getting directories is not supported for this client")
}

// GetFile downloads the give URL into the given path. The URL must
// reference a single file. If possible, the Getter should check if
// the remote end contains the same file and no-op this operation.
func (g *gsGetter) GetFile(dst string, u *url.URL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	client, err := g.client(ctx)
	if err != nil {
		return err
	}

	reader, err := g.objectAt(client, u).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("couldn't open object at %q: %v", u.String(), err)
	}

	return stream.Copy(stream.FromReadCloser(reader), stream.ToFile(dst))
}

func (gsGetter) objectAt(client *storage.Client, u *url.URL) *storage.ObjectHandle {
	return client.Bucket(u.Hostname()).Object(strings.TrimPrefix(u.Path, "/"))
}

func (gsGetter) client(ctx context.Context) (*storage.Client, error) {
	creds, err := google.CredentialsFromJSON(ctx, []byte(utils.GetServiceAccountJson()), storage.ScopeReadOnly)
	if err != nil {
		return nil, errors.New("couldn't get JSON credentials from the enviornment")
	}

	client, err := storage.NewClient(ctx, option.WithCredentials(creds), option.WithUserAgent(utils.CustomUserAgent))
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to Cloud Storage: %v", err)
	}
	return client, nil
}
