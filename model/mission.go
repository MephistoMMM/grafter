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
//
// Package model defined model to store mission datas to file system.
package model

// Mission represents a mission of grafting.
type Mission struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
	Name string `yaml:"name"`
}

// MissionStore store all registered Missions
type MissionStore struct {
	Version  string    `yaml:"version"`
	Missions []Mission `yaml:"missions"`
}

// NewMissionStore create and init a new MissionStore
func NewMissionStore(storePath string) (*MissionStore, error) {
	ms := &MissionStore{}
	err := ms.Init(storePath)
	return ms, err
}

// Init load mission store from yaml file
func (ms *MissionStore) Init(storePath string) error {
	ms.Version = "1.0.0"
	ms.Missions = []Mission{}
	return nil
}

// Add add mission to Missions field
func (ms *MissionStore) Add(mission Mission) bool {
	ms.Missions = append(ms.Missions, mission)
	return true
}
