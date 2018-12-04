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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBufferCallback_Write(t *testing.T) {
	writer := bufferCallback{}
	writer.Write([]byte("hello, world!"))

	if !reflect.DeepEqual(writer.Bytes(), []byte("hello, world!")) {
		t.Fatalf("write didn't append to the buffer")
	}
}

func TestBufferCallback_Close(t *testing.T) {
	testErr := errors.New("test error")
	calledBack := false
	writer := bufferCallback{
		closeCallback: func(b *bytes.Buffer) error {
			calledBack = true
			return testErr
		},
	}

	closeErr := writer.Close()

	if !calledBack {
		t.Errorf("Callback not called on close!")
	}

	if testErr != closeErr {
		t.Fatalf("expected err %v got %v", testErr, closeErr)
	}
}

func TestCopy(t *testing.T) {
	cases := map[string]struct {
		src      Source
		dest     Dest
		expected error
	}{
		"Source Init Err": {
			src:      FromError(errors.New("srcerr")),
			dest:     ToDiscard(),
			expected: errors.New("copy couldn't open source: srcerr"),
		},
		"Dest Init Err": {
			src:      FromString(""),
			dest:     ToError(errors.New("desterr")),
			expected: errors.New("copy couldn't open destination: desterr"),
		},

		"Dest Close Err": {
			src: FromString(""),
			dest: func() (io.WriteCloser, error) {
				return errWriteCloser{ioutil.Discard, errors.New("closerr")}, nil
			},
			expected: errors.New("copy couldn't close streams: 1 error(s) occurred: closerr"),
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			err := Copy(tc.src, tc.dest)
			if err.Error() != tc.expected.Error() {
				t.Fatalf("expected error: '%v' got '%v'", tc.expected, err)
			}
		})
	}
}

func ExampleYaml() {
	type Test struct {
		Str string `yaml:"s"`
		Num int    `yaml:"i"`
	}

	a := Test{"foo", 42}
	b := Test{}

	err := Copy(FromYaml(a), ToYaml(&b))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", b)

	// Output: stream.Test{Str:"foo", Num:42}
}

func ExampleFile() {
	td, err := ioutil.TempDir("", "test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(td)

	Copy(FromString("hello\nworld"), ToFile(td, "parent", "other", "testing.txt"))

	fullpath := filepath.Join(td, "parent", "other", "testing.txt")
	Copy(FromFile(fullpath), ToWriter(os.Stdout))

	// Output: hello
	// world
}

func ExampleMultiCloser() {
	var mc MultiCloser

	mc.Add(ioutil.NopCloser(nil))

	fmt.Println("closed:", mc.Close())

	// Output: closed: <nil>
}

func ExampleMultiCloser_Error() {
	var mc MultiCloser

	closer := errWriteCloser{nil, errors.New("example close error")}
	mc.Add(closer)

	fmt.Println("error closed:", mc.Close())
	fmt.Println("call after closed:", mc.Close())

	// Output: error closed: 1 error(s) occurred: example close error
	// call after closed: <nil>
}
