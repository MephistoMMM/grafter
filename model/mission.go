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

import (
	"github.com/MephistoMMM/grafter/util"
	"github.com/MephistoMMM/grafter/version"
	yaml "gopkg.in/yaml.v2"
)

// Mission represents a mission of grafting.
type Mission struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
	Name string `yaml:"name"`
}

// MissionStore store all registered Missions
type MissionStore struct {
	path     string
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
	ms.Version = version.Version
	ms.path = storePath

	return ms.Load(ms.path)
}

// Add add mission to Missions field
func (ms *MissionStore) Add(mission Mission) bool {
	// return false while mission already exists
	for _, m := range ms.Missions {
		if m.Name == mission.Name {
			return false
		}
	}

	ms.Missions = append(ms.Missions, mission)
	return true
}

// Path ...
func (ms *MissionStore) Path() string {
	return ms.path
}

// Load read mission data from path
func (ms *MissionStore) Load(path string) error {
	if util.IsNotExist(path) {
		ms.Missions = []Mission{}
		return nil
	}

	d, err := util.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(d, ms)
	if err != nil {
		return err
	}

	return nil
}

// Store write mission store data into path
func (ms *MissionStore) Store(path string) error {
	d, err := yaml.Marshal(ms)
	if err != nil {
		return err
	}

	// write mission datas to file system
	if err = util.WriteFile(path, d); err != nil {
		return err
	}

	return nil
}
