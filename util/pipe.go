// Copyright © 2018 Mephis Pheies <mephistommm@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package util

import (
	"os"
	"path/filepath"
)

// Item is the subject passed in pipe by channel
type Item struct {
	Err error

	Path string
}

// Walker walk the directory 'dir' to check each pathes of file undet it,
// then put path of valid file into pipe channel.
type Walker struct {
	dir     string
	checker IgnoreSupport

	pipe chan *Item
}

// NewWalker read and collect files recursively under the dir Directory
func NewWalker(dir string, checker IgnoreSupport, cap int) *Walker {
	return &Walker{
		dir:     dir,
		checker: checker,
		pipe:    make(chan *Item, cap),
	}
}

// Dir return dir
func (w Walker) Dir() string {
	return w.dir
}

// Pipe return single direction channel 'pipe'
func (w Walker) Pipe() <-chan *Item {
	return w.pipe
}

// Walk read and collect files recursively under the dir Directory
func (w *Walker) Walk() (err error) {
	err = filepath.Walk(w.dir,
		func(path string, info os.FileInfo, er error) error {
			if er != nil {
				// send error
				w.pipe <- &Item{
					Err: er,
				}
				return er
			}
			if !Check(w.checker, path, info) {
				// send valid path
				w.pipe <- &Item{
					Path: path,
				}
				return nil
			}

			if info.IsDir() {
				return filepath.SkipDir
			}

			return nil
		})

	close(w.pipe)
	return
}
