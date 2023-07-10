// Copyright © 2023 OpenIM SDK. All rights reserved.
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

package file

type PutFileCallback interface {
	Open(size int64)
	HashProgress(current, total int64)
	HashComplete(hash string, total int64)
	PutStart(current, total int64)
	PutProgress(save int64, current, total int64)
	PutComplete(total int64, putType int)
}

type emptyCallback struct{}

func (e emptyCallback) Open(size int64) {}

func (e emptyCallback) HashProgress(current, total int64) {}

func (e emptyCallback) HashComplete(hash string, total int64) {}

func (e emptyCallback) PutStart(current, total int64) {}

func (e emptyCallback) PutProgress(save int64, current, total int64) {}

func (e emptyCallback) PutComplete(total int64, putType int) {}
