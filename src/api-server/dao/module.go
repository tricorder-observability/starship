// Copyright (C) 2023  Tricorder Observability
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package dao

import (
	"fmt"

	"github.com/tricorder/src/utils/sqlite"
)

// ModuleGORM sqlite gorm storage and response object
type ModuleGORM struct {
	ID                 string `gorm:"'id' primarykey" json:"id,omitempty"`
	Name               string `gorm:"name" json:"name,omitempty"`
	DesiredState       int    `gorm:"desired_state" json:"desired_state,omitempty"`
	CreateTime         string `gorm:"create_time" json:"create_time,omitempty"`
	Ebpf               string `gorm:"ebpf" json:"ebpf,omitempty"`
	EbpfFmt            int    `gorm:"ebpf_fmt" json:"ebpf_fmt,omitempty"`
	EbpfLang           int    `gorm:"ebpf_lang" json:"ebpf_lang,omitempty"`
	EbpfPerfBufferName string `gorm:"ebpf_perf_name" json:"ebpf_perf_name,omitempty"`
	EbpfProbes         string `gorm:"ebpf_probes" json:"ebpf_probes,omitempty"`
	// wasm store the whole wasm file content
	Wasm       []byte `gorm:"wasm" json:"wasm,omitempty"`
	SchemaName string `gorm:"schema_name" json:"schema_name,omitempty"`
	SchemaAttr string `gorm:"schema_attr" json:"schema_attr,omitempty"`
	Fn         string `gorm:"fn" json:"fn,omitempty"`
	WasmFmt    int    `gorm:"wasm_fmt" json:"wasm_fmt,omitempty"`
	WasmLang   int    `gorm:"wasm_lang" json:"wasm_lang,omitempty"`
}

func (ModuleGORM) TableName() string {
	return "module"
}

type ModuleDao struct {
	Client *sqlite.ORM
}

func (g *ModuleDao) SaveModule(mod *ModuleGORM) error {
	result := g.Client.Engine.Create(mod)
	return result.Error
}

func (g *ModuleDao) UpdateByID(mod *ModuleGORM) error {
	if len(mod.ID) == 0 {
		return fmt.Errorf("code is 0")
	}

	result := g.Client.Engine.Model(&ModuleGORM{ID: mod.ID}).Updates(mod)
	return result.Error
}

func (g *ModuleDao) UpdateStatusByID(id string, desiredState int) error {
	result := g.Client.Engine.Model(&ModuleGORM{ID: id}).
		Select("desired_state").Updates(ModuleGORM{DesiredState: desiredState})
	return result.Error
}

func (g *ModuleDao) DeleteByID(id string) error {
	result := g.Client.Engine.Delete(&ModuleGORM{ID: id})
	return result.Error
}

func (g *ModuleDao) List(query ...string) ([]ModuleGORM, error) {
	var modeList []ModuleGORM
	if len(query) == 0 {
		query = []string{"id", "name", "desired_state", "create_time", "schema_attr", "fn", "ebpf"}
	}
	result := g.Client.Engine.
		Select(query).Where("name is not null and name != '' ").
		Order("create_time desc").
		Find(&modeList)
	if result.Error != nil {
		return nil, fmt.Errorf("query code list error:%v", result.Error)
	}
	return modeList, nil
}

func (g *ModuleDao) ListByStatus(desiredState int) ([]ModuleGORM, error) {
	var moduleList []ModuleGORM
	result := g.Client.Engine.Where(&ModuleGORM{DesiredState: desiredState}).Order("create_time desc").Find(&moduleList)
	if result.Error != nil {
		return make([]ModuleGORM, 0), fmt.Errorf("query code list by desiredStatus error:%v", result.Error)
	}
	return moduleList, nil
}

func (g *ModuleDao) QueryByName(name string) (*ModuleGORM, error) {
	code := &ModuleGORM{}
	result := g.Client.Engine.Where(&ModuleGORM{Name: name}).First(code)
	if result.Error != nil {
		return nil, fmt.Errorf("query code by name error:%v", result.Error)
	}
	return code, nil
}

func (g *ModuleDao) QueryByID(id string) (*ModuleGORM, error) {
	code := &ModuleGORM{}
	result := g.Client.Engine.Where(&ModuleGORM{ID: id}).First(code)
	if result.Error != nil {
		return nil, fmt.Errorf("query code by id error:%v", result.Error)
	}
	return code, nil
}
