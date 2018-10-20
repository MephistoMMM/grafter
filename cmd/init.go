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
	"log"
	"path/filepath"

	"github.com/MephistoMMM/grafter/model"
	"github.com/MephistoMMM/grafter/util"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init <mission_name> <SRC> <DEST>",
	Short: "Register a grafting mission.",
	Long:  `This command register a mission to associate SRC project with DEST project (SRC and DEST are both the path of project). It will write register information to $HOME/.grafter/missions, then compare all files in both projects to records temporary data to $HOME/.grafter/gm_$SRC@$DEST, which used in 'graft' command.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(3)(cmd, args); err != nil {
			return err
		}
		src, dest := args[1], args[2]
		if !util.IsDir(src) {
			return fmt.Errorf("src directory is not exist: %s", src)
		}
		if !util.IsDir(dest) {
			return fmt.Errorf("dest directory is not exist: %s", dest)
		}
		return nil
	},
	Run: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	name := args[0]
	srcDir, _ := filepath.Abs(args[1])
	destDir, _ := filepath.Abs(args[2])

	ok := Store.Add(model.Mission{
		Src:  srcDir,
		Dest: destDir,
		Name: name,
	})
	if !ok {
		log.Fatalf("Mission %s already exists.", name)
	}
}
