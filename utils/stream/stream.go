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

package stream

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
	multierror "github.com/hashicorp/go-multierror"
	yaml "gopkg.in/yaml.v2"
)

type Source func() (io.ReadCloser, error)
type Dest func() (io.WriteCloser, error)

// Copy copies data from a source stream to a destination stream.
func Copy(src Source, dest Dest) error {
	mc := MultiCloser{}
	defer mc.Close()

	readCloser, err := src()
	if err != nil {
		return fmt.Errorf("copy couldn't open source: %v", err)
	}
	mc.Add(readCloser)

	writeCloser, err := dest()
	if err != nil {
		return fmt.Errorf("copy couldn't open destination: %v", err)
	}
	mc.Add(writeCloser)

	if _, err := io.Copy(writeCloser, readCloser); err != nil {
		return fmt.Errorf("copy couldn't copy data: %v", err)
	}

	if err := mc.Close(); err != nil {
		return fmt.Errorf("copy couldn't close streams: %v", err)
	}

	return nil
}

// FromYaml converts the interface to a stream of Yaml.
func FromYaml(v interface{}) Source {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		return FromError(err)
	}

	return FromBytes(bytes)
}

// FromBytes streams the given bytes as a buffer.
func FromBytes(b []byte) Source {
	return FromReader(bytes.NewReader(b))
}

// FromBytes streams the given bytes as a buffer.
func FromString(s string) Source {
	return FromBytes([]byte(s))
}

// FromError returns a nil ReadCloser and the error passed when called.
func FromError(err error) Source {
	return func() (io.ReadCloser, error) {
		return nil, err
	}
}

// FromFile joins the segments of the path and reads from it.
func FromFile(path ...string) Source {
	return FromReadCloserError(os.Open(filepath.Join(path...)))
}

// FromReadCloserError reads the contents of the readcloser and takes ownership of closing it.
// If err is non-nil, it is returned as a source.
func FromReadCloserError(rc io.ReadCloser, err error) Source {
	return func() (io.ReadCloser, error) {
		return rc, err
	}
}

// FromReadCloser reads the contents of the readcloser and takes ownership of closing it.
func FromReadCloser(rc io.ReadCloser) Source {
	return FromReadCloserError(rc, nil)
}

// FromReader converts a Reader to a Source.
func FromReader(rc io.Reader) Source {
	return func() (io.ReadCloser, error) {
		return ioutil.NopCloser(rc), nil
	}
}

// ToFile concatenates the given path segments with filepath.Join, creates any
// parent directoreis if needed, and writes the file.
func ToFile(path ...string) Dest {
	return ToModeFile(0600, path...)
}

// ToModeFile is like ToFile, but sets the permissions on the created file.
func ToModeFile(mode os.FileMode, path ...string) Dest {
	return func() (io.WriteCloser, error) {
		outputPath := filepath.Join(path...)
		if err := os.MkdirAll(filepath.Dir(outputPath), 0700|os.ModeDir); err != nil {
			return nil, err
		}

		return os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	}
}

// ToBuffer buffers the contents of the stream and on Close() calls the callback
// returning its error.
func ToBuffer(closeCallback func(*bytes.Buffer) error) Dest {
	return func() (io.WriteCloser, error) {
		return &bufferCallback{closeCallback: closeCallback}, nil
	}
}

// ToYaml unmarshals the contents of the stream as YAML to the given struct.
func ToYaml(v interface{}) Dest {
	return ToBuffer(func(buf *bytes.Buffer) error {
		return yaml.Unmarshal(buf.Bytes(), v)
	})
}

// ToError returns an error when Dest is initialized.
func ToError(err error) Dest {
	return func() (io.WriteCloser, error) {
		return nil, err
	}
}

// ToDiscard discards any data written to it.
func ToDiscard() Dest {
	return ToWriter(ioutil.Discard)
}

// ToWriter forwards data to the given writer, this function WILL NOT close the
// underlying stream so it is safe to use with things like stdout.
func ToWriter(writer io.Writer) Dest {
	return ToWriteCloser(NopWriteCloser(writer))
}

// ToWriteCloser forwards data to the given WriteCloser which will be closed
// after the copy finishes.
func ToWriteCloser(w io.WriteCloser) Dest {
	return func() (io.WriteCloser, error) {
		return w, nil
	}
}

// bufferCallback buffers the results and on close calls the callback.
type bufferCallback struct {
	bytes.Buffer
	closeCallback func(*bytes.Buffer) error
}

// Close implements io.Closer.
func (b *bufferCallback) Close() error {
	return b.closeCallback(&b.Buffer)
}

// NopWriteCloser works like io.NopCloser, but for writers.
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return errWriteCloser{Writer: w, CloseErr: nil}
}

type errWriteCloser struct {
	io.Writer
	CloseErr error
}

// Close implements io.Closer.
func (w errWriteCloser) Close() error {
	return w.CloseErr
}

// MultiCloser calls Close() once on every added closer and returns
// a multierror of all errors encountered.
type MultiCloser struct {
	closers []io.Closer
	mu      sync.Mutex
}

// Add appends the closer to the internal list of closers.
func (mc *MultiCloser) Add(c io.Closer) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.closers = append(mc.closers, c)
}

// Close calls close on all closers, discarding them from the list.
// Modifying MultiCloser in your Close method will create a deadlock.
// At the end of this function, the close list will be empty and any encountered
// errors will be returned.
func (mc *MultiCloser) Close() error {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	result := &multierror.Error{ErrorFormat: utils.SingleLineErrorFormatter}

	for _, closer := range mc.closers {
		if err := closer.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	mc.closers = []io.Closer{}
	return result.ErrorOrNil()
}
