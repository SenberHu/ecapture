// Copyright 2022 CFC4N <cfc4n.cs@gmail.com>. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"debug/elf"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	GoDashReadFunc = "read"
)

// DashConfig
type DashConfig struct {
	BaseConfig
	Dashpath         string `json:"dashpath"` //dash的文件路径
	Address          uint64
	ErrNo            int
	ElfType          uint8
	goElfArch        string    //
	goElf            *elf.File //
	IsPieBuildMode   bool
	ReadlineFuncName string
}

func NewDashConfig() *DashConfig {
	config := &DashConfig{}
	config.PerCpuMapSize = DefaultMapSizePerCpu
	return config
}

func (dc *DashConfig) findRetOffsets(symbolName string) ([]int, error) {
	var err error
	var allSymbs []elf.Symbol

	goSymbs, _ := dc.goElf.Symbols()
	if len(goSymbs) > 0 {
		allSymbs = append(allSymbs, goSymbs...)
	}
	goDynamicSymbs, _ := dc.goElf.DynamicSymbols()
	if len(goDynamicSymbs) > 0 {
		allSymbs = append(allSymbs, goDynamicSymbs...)
	}

	if len(allSymbs) == 0 {
		return nil, ErrorSymbolEmpty
	}
	var nameList string
	var found bool
	var symbol elf.Symbol
	for _, s := range allSymbs {
		// if s.Name == symbolName {
		// 	symbol = s
		// 	found = true

		// 	break
		// }
		nameList += fmt.Sprintf("name%s address0x%x|", s.Name, s.Value)
	}
	return nil, fmt.Errorf("symbol not found: %s\n%s", symbolName, nameList)

	if !found {
		return nil, ErrorSymbolNotFound
	}

	section := dc.goElf.Sections[symbol.Section]

	var elfText []byte
	elfText, err = section.Data()
	if err != nil {
		return nil, err
	}

	start := symbol.Value - section.Addr
	end := start + symbol.Size

	var offsets []int
	var instHex []byte
	instHex = elfText[start:end]
	offsets, _ = DecodeInstruction(instHex)
	if len(offsets) == 0 {
		return offsets, ErrorNoRetFound
	}

	address := symbol.Value
	for _, prog := range dc.goElf.Progs {
		// Skip uninteresting segments.
		if prog.Type != elf.PT_LOAD || (prog.Flags&elf.PF_X) == 0 {
			continue
		}

		if prog.Vaddr <= symbol.Value && symbol.Value < (prog.Vaddr+prog.Memsz) {
			// stackoverflow.com/a/40249502
			address = symbol.Value - prog.Vaddr + prog.Off
			break
		}
	}
	for i, offset := range offsets {
		offsets[i] = int(address) + offset
	}
	return offsets, nil
}

func (dc *DashConfig) Check() error {
	if dc.Dashpath == "" || len(strings.TrimSpace(dc.Dashpath)) == 0 {
		dc.Dashpath = "/bin/dash"
	}

	_, err := os.Stat(dc.Dashpath)
	if err != nil {
		fmt.Printf("111")
		return err
	}
	dc.ElfType = ElfTypeBin

	//如果配置 funcname ，则使用用户指定的函数名
	if dc.ReadlineFuncName == "" || len(strings.TrimSpace(dc.ReadlineFuncName)) == 0 {
		dc.ReadlineFuncName = "read@plt"
	}

	// var goElf *elf.File
	// goElf, err = elf.Open(dc.Dashpath)
	// if err != nil {
	// 	fmt.Printf("444")
	// 	return err
	// }

	// var goElfArch string
	// switch goElf.FileHeader.Machine.String() {
	// case elf.EM_AARCH64.String():
	// 	goElfArch = "arm64"
	// case elf.EM_X86_64.String():
	// 	goElfArch = "amd64"
	// default:
	// 	goElfArch = "unsupport_arch"
	// }

	// if goElfArch != runtime.GOARCH {
	// 	err = fmt.Errorf("Go Application not match, want:%s, have:%s", runtime.GOARCH, goElfArch)
	// 	return err
	// }
	// switch goElfArch {
	// case "amd64":
	// case "arm64":
	// default:
	// 	return fmt.Errorf("unsupport CPU arch :%s", goElfArch)
	// }

	// if dc.Address != 0 {
	// 	return nil
	// }

	// dc.goElfArch = goElfArch
	// dc.goElf = goElf
	// // If built with PIE and stripped, gopclntab is
	// // unlabeled and nested under .data.rel.ro.

	// var addr []int
	// addr, err = dc.findRetOffsets(GoDashReadFunc)
	// if err != nil {
	// 	return fmt.Errorf("%s symbol address error:%s", GoDashReadFunc, err.Error())
	// }
	// var addrlist string
	// for _, a := range addr {
	// 	addrlist += fmt.Sprintf("%d |", a)
	// }

	//return fmt.Errorf("addrt is %d", addrlist)
	return nil
}

func (dc *DashConfig) checkElf() error {
	//如果配置 dash的地址，且存在，则直接返回
	if dc.Dashpath != "" || len(strings.TrimSpace(dc.Dashpath)) > 0 {
		_, e := os.Stat(dc.Dashpath)
		if e != nil {
			return e
		}
		dc.ElfType = ElfTypeBin
		return nil
	}

	return nil
}

func (dc *DashConfig) Bytes() []byte {
	b, e := json.Marshal(dc)
	if e != nil {
		return []byte{}
	}
	return b
}
