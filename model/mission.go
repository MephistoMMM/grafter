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
	"fmt"

	"github.com/MephistoMMM/grafter/util"
	"github.com/MephistoMMM/grafter/version"
	yaml "gopkg.in/yaml.v2"
)

// Mission represents a mission of grafting.
type Mission struct {
	Src    string   `yaml:"src"`
	Dest   string   `yaml:"dest"`
	Name   string   `yaml:"name"`
	Ignore []string `yaml:"ignore"`
}

// String return string value of Mission data
func (m *Mission) String() string {
	return fmt.Sprintf("%s:\n\tsrc: %s\n\tdest: %s\n\tignore: %s\n",
		m.Name, m.Src, m.Dest, m.Ignore)
}

// AddIgnore append a new regex string to Ignore field, if it
// has not been existed.
func (m *Mission) AddIgnore(reStr string) {
	for _, i := range m.Ignore {
		if i == reStr {
			return
		}
	}

	m.Ignore = append(m.Ignore, reStr)
}

// RemoveIgnore delete a regex string from Ignore field. if it
// has not been existed, do nothing.
func (m *Mission) RemoveIgnore(index int64) {
	if index == 0 {
		m.Ignore = m.Ignore[1:]
	} else if index > 0 && index < int64(len(m.Ignore)-1) {
		m.Ignore = append(m.Ignore[0:index], m.Ignore[index+1:]...)
	} else if index == int64(len(m.Ignore)-1) {
		m.Ignore = m.Ignore[0 : len(m.Ignore)-1]
	}
}

// MissionStore store all registered Missions
type MissionStore struct {
	path     string
	modified bool
	Version  string    `yaml:"version"`
	Missions []Mission `yaml:"missions"`
}

// NewMissionStore create and init a new MissionStore
func NewMissionStore(storePath string) (*MissionStore, error) {
	ms := &MissionStore{}
	err := ms.Init(storePath)
	return ms, err
}

// Modified set the value of modified field
func (ms *MissionStore) Modified(t bool) {
	ms.modified = t
}

// Init load mission store from yaml file
func (ms *MissionStore) Init(storePath string) error {
	ms.Version = version.Version
	ms.path = storePath
	ms.modified = false

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
	ms.modified = true
	return true
}

// Get mission by mission name
func (ms *MissionStore) Get(mission string) *Mission {
	// return false while mission already exists
	for i, m := range ms.Missions {
		if m.Name == mission {
			return &ms.Missions[i]
		}
	}

	return nil
}

// Remove remove mission to Missions field
func (ms *MissionStore) Remove(mission string) bool {
	// return false while mission already exists
	for i, m := range ms.Missions {
		if m.Name == mission {
			ms.Missions = append(ms.Missions[:i], ms.Missions[i+1:]...)
			ms.modified = true
			return true
		}
	}

	return false
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
	// do nothing if destination path is same and mission store is not
	// modified
	if path == ms.path && !ms.modified {
		return nil
	}

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
