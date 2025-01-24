// Copyright © 2023 OpenIM open source community. All rights reserved.
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

//go:build js && wasm

package file

import (
	"bufio"
	"errors"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
	"io"
	"syscall/js"
)

const readBufferSize = 1024 * 1024 * 5 // 5mb

func Open(req *UploadFileReq) (ReadFile, error) {
	file := newJsCallFile(req.Uuid)
	size, err := file.Open()
	if err != nil {
		return nil, err
	}
	jf := &jsFile{
		size: size,
		file: file,
	}
	jf.resetReaderBuffer()
	return jf, nil
}

type jsFile struct {
	size   int64
	file   *jsCallFile
	whence int
	reader io.Reader
}

func (j *jsFile) resetReaderBuffer() {
	j.reader = bufio.NewReaderSize(&reader{fn: j.read}, readBufferSize)
}

func (j *jsFile) read(p []byte) (n int, err error) {
	length := len(p)
	if length == 0 {
		return 0, errors.New("read buffer is empty")
	}
	if j.whence >= int(j.size) {
		return 0, io.EOF
	}
	if j.whence+length > int(j.size) {
		length = int(j.size) - j.whence
	}
	data, err := j.file.Read(int64(j.whence), int64(length))
	if err != nil {
		return 0, err
	}
	if len(data) > len(p) {
		return 0, errors.New("js read data > length")
	}
	j.whence += len(data)
	copy(p, data)
	return len(data), nil
}

func (j *jsFile) Read(p []byte) (n int, err error) {
	return j.reader.Read(p)
}

func (j *jsFile) Close() error {
	return j.file.Close()
}

func (j *jsFile) Size() int64 {
	return j.size
}

func (j *jsFile) StartSeek(whence int) error {
	if whence < 0 || whence > int(j.size) {
		return errors.New("seek whence is out of range")
	}
	j.whence = whence
	j.resetReaderBuffer()
	return nil
}

type reader struct {
	fn func(p []byte) (n int, err error)
}

func (r *reader) Read(p []byte) (n int, err error) {
	return r.fn(p)
}

type jsCallFile struct {
	uuid string
}

func newJsCallFile(uuid string) *jsCallFile {
	return &jsCallFile{uuid: uuid}
}

func (j *jsCallFile) Open() (int64, error) {
	return WasmOpen(j.uuid)
}

func (j *jsCallFile) Read(offset int64, length int64) ([]byte, error) {
	return WasmRead(j.uuid, offset, length)
}

func (j *jsCallFile) Close() error {
	return WasmClose(j.uuid)
}

func WasmOpen(uuid string) (int64, error) {
	result, err := exec.Exec(uuid)
	if err != nil {
		return 0, err
	}
	if v, ok := result.(float64); ok {
		size := int64(v)
		if size < 0 {
			return 0, errors.New("file size < 0")
		}
		return size, nil
	}
	return 0, exec.ErrType
}
func WasmRead(uuid string, offset int64, length int64) ([]byte, error) {
	result, err := exec.Exec(uuid, offset, length)
	if err != nil {
		return nil, err
	} else {
		if v, ok := result.(js.Value); ok {
			return exec.ExtractArrayBuffer(v), nil
		} else {
			return nil, exec.ErrType
		}
	}
}
func WasmClose(uuid string) error {
	_, err := exec.Exec(uuid)
	return err
}
