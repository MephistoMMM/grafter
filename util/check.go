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
	"regexp"

	"github.com/MephistoMMM/grafter/gitignore"
)

// Walk is designed as `chain of repositories` mode. It includes three parts,
// 1.a abstract Support class implemented by IgnoreSupport interface and
// BaseSupport struct , 2.some concrete Support, 3.Check function with IgnoreSupport
// parameter calling each IsIgnore method of each IgnoreSupport in chain of repositories.

// IgnoreSupport is a part of abstract Support.
type IgnoreSupport interface {
	SetNext(IgnoreSupport) IgnoreSupport
	SetNexts([]IgnoreSupport) IgnoreSupport
	Next() IgnoreSupport
	String() string

	IsIgnore(path string, info os.FileInfo) (bool, error)
	Done(path string, info os.FileInfo)
	Fail(path string, info os.FileInfo)
}

// BaseSupport implements almost methods of IgnoreSupport interface except
// `IsIgnore`. It should be embeded into a concrete IgnoreSupport.
type BaseSupport struct {
	next IgnoreSupport
	name string
}

// SetName set n to name
// This method should only be used when concrete IgnoreSupport inits, so it
// is not a part of IgnoreSupport interface.
func (bs *BaseSupport) SetName(n string) {
	bs.name = n
}

// SetNext assign another IgnoreSupport to inner next field. It is used to
// construct a chain of IgnoreSupports.
func (bs *BaseSupport) SetNext(n IgnoreSupport) IgnoreSupport {
	bs.next = n
	return n
}

// SetNexts assign a IgnoreSupport link list to inner next field.
func (bs *BaseSupport) SetNexts(ns []IgnoreSupport) IgnoreSupport {
	if len(ns) < 1 {
		return nil
	}
	innerBs := bs.SetNext(ns[0])
	ns = ns[1:]
	for _, n := range ns {
		innerBs = innerBs.SetNext(n)
	}
	return innerBs
}

// Next return next.
func (bs *BaseSupport) Next() IgnoreSupport {
	return bs.next
}

// String describe the chain of IgnoreSupport
func (bs *BaseSupport) String() string {
	if bs.Next() == nil {
		return bs.name
	}
	return bs.name + " | " + bs.Next().String()
}

// Done does nothing but implement IgnoreSupport interface
func (bs *BaseSupport) Done(path string, info os.FileInfo) {
	Logger.Debugf("%s ignored by %s\n", path, bs.name)
}

// Fail does nothing but implement IgnoreSupport interface
func (bs *BaseSupport) Fail(path string, info os.FileInfo) {
}

// Check calls each IsIgnore method of each IgnoreSupport in chain of repositories.
func Check(checker IgnoreSupport, path string, info os.FileInfo) bool {
	if checker == nil {
		return true
	}

	for checker != nil {
		result, err := checker.IsIgnore(path, info)
		if err != nil {
			checker.Fail(path, info)
		}

		if result {
			checker.Done(path, info)
			return true
		}

		checker = checker.Next()
	}

	return false
}

type IgnoreRegexpMatchSupport struct {
	BaseSupport

	pattern *regexp.Regexp
}

func NewMultiIgnoreRegexpMatchSupports(exprs []string) ([]IgnoreSupport, error) {
	iss := make([]IgnoreSupport, 0, len(exprs))
	for _, exprs := range exprs {
		is, err := NewIgnoreRegexpMatchSupport(exprs)
		if err != nil {
			return nil, err
		}

		iss = append(iss, is)
	}
	return iss, nil
}

func NewIgnoreRegexpMatchSupport(expr string) (IgnoreSupport, error) {
	pattern, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}

	is := &IgnoreRegexpMatchSupport{
		pattern: pattern,
	}
	is.SetName("IgnoreRegexpMatchSupport")
	return is, nil
}

// IsIgnore ...
func (irms *IgnoreRegexpMatchSupport) IsIgnore(path string, info os.FileInfo) (bool, error) {
	if irms.pattern.MatchString(path) {
		return true, nil
	}
	return false, nil
}

// IgnoreSpecialMadeSupport ignore special type fies.
type IgnoreSpecialMadeSupport struct {
	BaseSupport

	modeMask os.FileMode
}

func NewIgnoreUnregularSupport() (IgnoreSupport, error) {
	is := &IgnoreSpecialMadeSupport{
		modeMask: os.ModeSymlink | os.ModeNamedPipe | os.ModeSocket | os.ModeDevice | os.ModeIrregular,
	}
	is.SetName("IgnoreUnregularSupport")
	return is, nil
}

// IsIgnore ...
func (isms *IgnoreSpecialMadeSupport) IsIgnore(path string, info os.FileInfo) (bool, error) {
	if isms.modeMask&info.Mode() != 0 {
		return true, nil
	}

	return false, nil
}

// IgnoreDotSupport ignore all dot fies.
type IgnoreDotSupport struct {
	BaseSupport
}

func NewIgnoreDotSupport() (IgnoreSupport, error) {
	is := &IgnoreDotSupport{}
	is.SetName("IgnoreDotSupport")
	return is, nil
}

// IsIgnore ...
func (ids *IgnoreDotSupport) IsIgnore(path string, info os.FileInfo) (bool, error) {
	if filepath.Base(path)[0] == '.' {
		return true, nil
	}

	return false, nil
}

type GitIgnoreSupport struct {
	BaseSupport

	checker *gitignore.Checker
}

func NewGitIgnoreSupport(path string) (IgnoreSupport, error) {
	checker := gitignore.NewChecker()
	if err := checker.LoadBasePath(path); err != nil {
		return nil, err
	}

	is := &GitIgnoreSupport{
		checker: checker,
	}
	is.SetName("GitIgnoreSupport")
	return is, nil
}

// IsIgnore ...
func (gis *GitIgnoreSupport) IsIgnore(path string, info os.FileInfo) (bool, error) {
	return gis.checker.Check(path, info), nil
}
