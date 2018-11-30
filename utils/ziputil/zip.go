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

package ziputil

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/GoogleCloudPlatform/gcp-service-broker/utils/stream"
)

// List writes a ls -la style listing of the zipfile to the given writer.
func List(z *zip.Reader, w io.Writer) {
	sw := tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.StripEscape)
	fmt.Fprintln(sw, "MODE\tSIZE\tNAME")
	for _, fd := range z.File {
		fmt.Fprintf(sw, "%s\t%d\t%s\n", fd.Mode().String(), fd.UncompressedSize, fd.Name)
	}
	sw.Flush()
}

// Joins a path for use in a zip file.
func Join(path ...string) string {
	return strings.Join(path, "/")
}

func Clean(path ...string) string {
	joined := filepath.ToSlash(Join(path...))
	slashStrip := strings.TrimPrefix(joined, "/")
	dotStrip := strings.TrimPrefix(slashStrip, "./")
	return dotStrip
}

// Find returns a pointer to the file at the given path if it exists, otherwise
// nil.
func Find(z *zip.Reader, path ...string) *zip.File {
	name := Join(path...)
	for _, f := range z.File {
		if f.Name == name {
			return f
		}
	}

	return nil
}

// Opens the file at the given path if possible, otherwise returns an error.
func Open(z *zip.Reader, path ...string) (io.ReadCloser, error) {
	f := Find(z, path...)
	if f == nil {
		fmt.Errorf("no such file: %q", Join(path...))
	}

	return f.Open()
}

// Extracts the contents of the zipDirectory to the given OS osDirectory.
func Extract(z *zip.Reader, zipDirectory, osDirectory string) error {
	for _, fd := range z.File {
		if fd.UncompressedSize == 0 { // skip directories
			continue
		}

		if !strings.HasPrefix(fd.Name, zipDirectory) {
			continue
		}

		src := stream.FromReadCloserError(fd.Open())

		newName := strings.TrimPrefix(fd.Name, zipDirectory)
		destPath := filepath.Join(osDirectory, filepath.FromSlash(newName))
		dest := stream.ToModeFile(fd.Mode(), destPath)

		if err := stream.Copy(src, dest); err != nil {
			return fmt.Errorf("couldn't extract file %q: %v", fd.Name, err)
		}
	}

	return nil
}

// Unarchive opens the specified file and extracts all of its contents to the
// destination.
func Unarchive(zipFile, destination string) error {
	rdc, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("couldn't open archive %q: %v", zipFile, err)
	}
	defer rdc.Close()
	return Extract(&rdc.Reader, "", destination)
}

// Archive creates a zip from the contents of the sourceFolder at the
// destinationZip location.
func Archive(sourceFolder, destinationZip string) error {
	fd, err := os.Create(destinationZip)
	if err != nil {
		return fmt.Errorf("couldn't create archive %q: %v", destinationZip, err)
	}

	defer fd.Close()

	w := zip.NewWriter(fd)
	defer w.Close()

	err = filepath.Walk(sourceFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = Clean(strings.TrimPrefix(path, sourceFolder))

		if info.IsDir() {
			w.CreateHeader(header)
		} else {
			fd, err := w.CreateHeader(header)
			if err != nil {
				return err
			}

			if err := stream.Copy(stream.FromFile(path), stream.ToWriter(fd)); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
