//    Copyright 2021 Olivier Mengué
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// package multidim initializes multidimensional slices.
//
// See https://golang.org/doc/effective_go#two_dimensional_slices
package multidim

import (
	"errors"
	"fmt"
	"reflect"
)

// Init allocates a multidimensional slice and initializes it.
//
// The size of each dimension to initialize must be given
// (but no dimension also works and initializes the single value *target).
//
// If build is not nil, it is used to initialize each cell.
// build is either (1) a direct value, (2) an anonymous function producing
// each cell value (func() cellT), or (3) an anonmymous function initializing
// each cell (func(&cellT)).
func Init(target interface{}, build interface{}, dimensions ...int) {
	nbDim := len(dimensions)
	targetV := reflect.ValueOf(&target).Elem().Elem()
	targetT := targetV.Type()
	if targetT.Kind() != reflect.Ptr || targetV.IsNil() {
		if len(dimensions) > 0 {
			panic(errors.New("target must be a pointer to a slice"))
		} else {
			panic(errors.New("target must be a pointer"))
		}
	}
	targetV = targetV.Elem()
	targetT = targetT.Elem()
	buildV := reflect.ValueOf(&build).Elem()
	if len(dimensions) == 0 {
		if build == nil {
			return
		}
		buildV = buildV.Elem()
		buildT := buildV.Type()
		if buildT.Kind() == reflect.Func && buildT.Name() == "" {
			// FIXME handle <func(&cellT)> and <func() cellT>
			if buildT.NumIn() == 1 {
				if buildT.NumOut() == 0 {
					buildV.Call([]reflect.Value{targetV.Addr()})
					return
				}
			} else if buildT.NumOut() == 1 {
				if buildT.NumIn() == 0 {
					targetV.Set(buildV.Call(nil)[0])
					return
				}
			}
			panic(fmt.Errorf("build function %T is not handled", build))
		}
		// fmt.Println(targetV.Type(), buildV.Type())
		targetV.Set(buildV)
		return
	}

	var size int = 1
	t := targetT
	types := make([]reflect.Type, nbDim)
	for d, dimSize := range dimensions {
		if t.Kind() != reflect.Slice {
			panic(fmt.Errorf("%s: dimension %d is not a slice", targetT, d+1))
		}
		if dimSize <= 0 {
			panic(fmt.Errorf("dimension %d: invalid size %d", d+1, dimSize))
		}
		types[d] = t
		size = size * dimSize
		t = t.Elem()
	}
	// fmt.Println(types)
	cellT := t
	cells := reflect.MakeSlice(types[len(types)-1], size, size)

	if build != nil {
		buildV = buildV.Elem()
		buildT := buildV.Type()
		if buildT.Kind() == reflect.Func && buildT.Name() == "" {
			numOut := buildT.NumOut()
			if buildT.NumIn() == 0 {
				if numOut == 1 {
					// Handle <func() cellT>
					for i := 0; i < size; i++ {
						cells.Index(i).Set(buildV.Call(nil)[0])
					}
				}
				goto ok
			} else if buildT.NumIn() == 1 && buildT.In(0) == reflect.PtrTo(cellT) && numOut == 0 {
				// Handle <func(&cellT)>
				args := make([]reflect.Value, 1)
				for i := 0; i < size; i++ {
					args[0] = cells.Index(i).Addr()
					buildV.Call(args)
				}
				goto ok
			}
			// FIXME handle init with <func(&cellT, dim1 int...)>
			// FIXME handle init with <func(dim1 int...) cellT>
			// make([]reflect.Value, nbDim)
			panic(fmt.Errorf("build function %T is not handled", build))
		ok:
		} else {
			for i := 0; i < size; i++ {
				cells.Index(i).Set(buildV)
			}
		}
	}

	nbCells := size
	for d := nbDim - 2; d >= 0; d-- {
		dimSize := dimensions[d+1]
		nbSlices := nbCells / dimSize
		slices := reflect.MakeSlice(types[d], nbSlices, nbSlices)
		offset := 0
		for i := 0; i < nbSlices; i++ {
			// slices[i] = cells[:dimSize:dimSize]
			slices.Index(i).Set(cells.Slice3(offset, offset+dimSize, offset+dimSize))
			offset += dimSize
		}
		cells = slices
		nbCells = nbSlices
	}

	targetV.Set(cells)
}
