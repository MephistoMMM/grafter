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

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/MephistoMMM/grafter/model"
	"github.com/MephistoMMM/grafter/util"
	"github.com/spf13/cobra"
)

// graftCmd represents the graft command
var graftCmd = &cobra.Command{
	Use:   "graft <mission_name>",
	Short: "Move files to DEST project from the SRC",
	Long: `Graft command just move all files except ignored files to DEST project from the SRC. More powerful functions will coming soon.
	Ignore configs should be add into the item of grafter configuration file.
`,
	Args: cobra.ExactArgs(1),
	Run:  graftRun,
}

func init() {
	rootCmd.AddCommand(graftCmd)
}

func graftRun(cmd *cobra.Command, args []string) {
	name := args[0]

	M := Store.Get(name)
	if M == nil {
		log.Fatalf("Mission %s doesn't exist.", name)
	}

	graft(M)
}

func combineIgnoreChain(M *model.Mission) util.IgnoreSupport {
	var checker, tail util.IgnoreSupport
	tail, err := util.NewIgnoreDotSupport()
	if err != nil {
		log.Fatal(err)
	}
	checker = tail

	ignoreUnregular, err := util.NewIgnoreUnregularSupport()
	if err != nil {
		log.Fatal(err)
	}
	tail = tail.SetNext(ignoreUnregular)

	// gitignore support ignores filepath matched patterns in .gitignore
	gitIgnore, err := util.NewGitIgnoreSupport(M.Src)
	if err != nil {
		log.Fatal(err)
	}
	tail = tail.SetNext(gitIgnore)

	// create regexp match supports from ignore field
	regexpMatches, err := util.NewMultiIgnoreRegexpMatchSupports(M.Ignore)
	if err != nil {
		log.Fatal(err)
	}
	tail = tail.SetNexts(regexpMatches)

	return checker
}

func graft(M *model.Mission) {
	log.Infof("Do Graft For %s", M.Name)

	checker := combineIgnoreChain(M)
	copyDifferent(M, checker)
	removeNotExist(M, checker)
}

func removeNotExist(M *model.Mission, checker util.IgnoreSupport) {
	walker := util.NewWalker(M.Dest, checker, 10)
	go func(w *util.Walker) {
		if err := w.Walk(); err != nil {
			log.Errorf("Walker[%s] error: %v.", w.Dir(), err)
		}
	}(walker)

	var wg sync.WaitGroup
	pipe := walker.Pipe()
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go doRemove(&wg, pipe, M)
	}

	wg.Wait()
}

func copyDifferent(M *model.Mission, checker util.IgnoreSupport) {
	walker := util.NewWalker(M.Src, checker, 10)
	go func(w *util.Walker) {
		if err := w.Walk(); err != nil {
			log.Errorf("Walker[%s] error: %v.", w.Dir(), err)
		}
	}(walker)

	var wg sync.WaitGroup
	pipe := walker.Pipe()
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go doCopy(&wg, pipe, M)
	}

	wg.Wait()
}

func doRemove(wg *sync.WaitGroup, pipe <-chan *util.Item, M *model.Mission) {
	for dest := range pipe {
		if dest.Err != nil {
			log.Infof("Receive Error From Pipe: %v.", dest.Err)
			continue
		}

		src := filepath.Join(M.Src, dest.Path[len(M.Dest):])
		if _, err := os.Stat(src); !os.IsNotExist(err) {
			continue
		}

		if err := os.Remove(dest.Path); err != nil {
			log.Errorf("Failed to remove %s: %s", dest.Path, err.Error())
		}

	}

	wg.Done()
}

func doCopy(wg *sync.WaitGroup, pipe <-chan *util.Item, M *model.Mission) {
	for source := range pipe {
		if source.Err != nil {
			log.Infof("Receive Error From Pipe: %v.", source.Err)
			continue
		}

		dest := filepath.Join(M.Dest, source.Path[len(M.Src):])
		isSame, err := compareFile(source.Path, dest)
		if err != nil {
			log.Error(err.Error())
		}

		if !isSame {
			util.CopyFile(source.Path, dest)
		}
	}

	wg.Done()
}

func compareFile(src, dest string) (bool, error) {
	srcFi, err := os.Lstat(src)
	if os.IsNotExist(err) {
		return false, fmt.Errorf("Source file %s doesn't exist!", src)
	}

	destFi, err := os.Lstat(dest)
	if os.IsNotExist(err) {
		return false, nil
	}

	if srcFi.Size() != destFi.Size() {
		return false, nil
	}

	// TODO : compare the content of two file
	return false, nil
}
