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
	"path/filepath"

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

func graft(M *model.Mission) {
	log.Infof("Do Graft For %s", M.Name)
	var checker util.IgnoreSupport
	ignoreDot, err := util.NewIgnoreDotSupport()
	if err != nil {
		log.Fatal(err)
	}

	ignoreUnregular, err := util.NewIgnoreUnregularSupport()
	if err != nil {
		log.Fatal(err)
	}

	// gitignore support ignores filepath matched patterns in .gitignore
	gitIgnore, err := util.NewGitIgnoreSupport(M.Src)
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln(M.Src)

	ignoreDot.SetNext(ignoreUnregular).SetNext(gitIgnore)
	checker = ignoreDot

	sources, _ := util.Walk(M.Src, checker)
	for _, source := range sources {
		dest := filepath.Join(M.Dest, source[len(M.Src):])
		// log.Debug(dest)
		util.CopyFile(source, dest)
	}
}
