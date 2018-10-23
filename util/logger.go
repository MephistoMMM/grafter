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
	"log"
	"os"
)

var Logger MixLogger

type MixLogger struct {
	out *log.Logger
	err *log.Logger
}

func init() {
	Logger.out = log.New(os.Stdout, "", 0)
	Logger.err = log.New(os.Stderr, "", 0)
}

func (ml *MixLogger) Print(v ...interface{}) {
	ml.out.Print(v...)
}

func (ml *MixLogger) Println(v ...interface{}) {
	ml.out.Println(v...)
}

func (ml *MixLogger) Printf(format string, v ...interface{}) {
	ml.out.Printf(format, v...)
}

func (ml *MixLogger) Fatal(v ...interface{}) {
	ml.err.Fatal(v...)
}

func (ml *MixLogger) Fatalln(v ...interface{}) {
	ml.err.Fatalln(v...)
}

func (ml *MixLogger) Fatalf(format string, v ...interface{}) {
	ml.err.Fatalf(format, v...)
}

func Print(v ...interface{}) {
	Logger.Print(v...)
}

func Println(v ...interface{}) {
	Logger.Println(v...)
}

func Printf(format string, v ...interface{}) {
	Logger.Printf(format, v...)
}

func Fatal(v ...interface{}) {
	Logger.Fatal(v...)
}

func Fatalln(v ...interface{}) {
	Logger.Fatalln(v...)
}

func Fatalf(format string, v ...interface{}) {
	Logger.Fatalf(format, v...)
}
