/*
Copyright 2021 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"go/types"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/packages"
)

const licenseFileHeader = `Copyright 2021 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.`

type sizegen struct {
	DebugTypes bool
	mod        *packages.Module
	sizes      types.Sizes
	codegen    map[string]*codeFile
	known      map[*types.Named]*typeState
}

type generatedCode struct {
	mod   *packages.Module
	files map[string]*codeFile
}

type codeFlag uint32

const (
	codeWithInterface = 1 << 0
	codeWithUnsafe    = 1 << 1
)

type codeImpl struct {
	name  string
	flags codeFlag
	code  jen.Code
}

type codeFile struct {
	pkg   string
	impls []codeImpl
}

type typeState struct {
	generated bool
	local     bool
	pod       bool // struct with only primitives
}

func newSizegen(mod *packages.Module, sizes types.Sizes) *sizegen {
	return &sizegen{
		DebugTypes: true,
		mod:        mod,
		sizes:      sizes,
		known:      make(map[*types.Named]*typeState),
		codegen:    make(map[string]*codeFile),
	}
}

func isPod(tt types.Type) bool {
	switch tt := tt.(type) {
	case *types.Struct:
		for i := 0; i < tt.NumFields(); i++ {
			if !isPod(tt.Field(i).Type()) {
				return false
			}
		}
		return true

	case *types.Basic:
		switch tt.Kind() {
		case types.String, types.UnsafePointer:
			return false
		}
		return true

	default:
		return false
	}
}

func (sizegen *sizegen) getKnownType(named *types.Named) *typeState {
	ts := sizegen.known[named]
	if ts == nil {
		local := strings.HasPrefix(named.Obj().Pkg().Path(), sizegen.mod.Path)
		ts = &typeState{
			local: local,
			pod:   isPod(named.Underlying()),
		}
		sizegen.known[named] = ts
	}
	return ts
}

func (sizegen *sizegen) generateType(pkg *types.Package, file *codeFile, named *types.Named) {
	ts := sizegen.getKnownType(named)
	if ts.generated {
		return
	}
	ts.generated = true

	switch tt := named.Underlying().(type) {
	case *types.Struct:
		if impl, flag := sizegen.sizeImplForStruct(named.Obj(), tt); impl != nil {
			file.impls = append(file.impls, codeImpl{
				code:  impl,
				name:  named.String(),
				flags: flag,
			})
		}
	case *types.Interface:
		findImplementations(pkg.Scope(), tt, func(tt types.Type) {
			if _, isStruct := tt.Underlying().(*types.Struct); isStruct {
				sizegen.generateType(pkg, file, tt.(*types.Named))
			}
		})
	default:
		// no-op
	}
}

func (sizegen *sizegen) generateKnownType(named *types.Named) {
	pkgInfo := named.Obj().Pkg()
	file := sizegen.codegen[pkgInfo.Path()]
	if file == nil {
		file = &codeFile{pkg: pkgInfo.Name()}
		sizegen.codegen[pkgInfo.Path()] = file
	}

	sizegen.generateType(pkgInfo, file, named)
}

func findImplementations(scope *types.Scope, iff *types.Interface, impl func(types.Type)) {
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		baseType := obj.Type()
		if types.Implements(baseType, iff) || types.Implements(types.NewPointer(baseType), iff) {
			impl(baseType)
		}
	}
}

func (sizegen *sizegen) generateKnownInterface(pkg *types.Package, iff *types.Interface) {
	findImplementations(pkg.Scope(), iff, func(tt types.Type) {
		if named, ok := tt.(*types.Named); ok {
			sizegen.generateKnownType(named)
		}
	})
}

func (sizegen *sizegen) finalize() {
	code := sizegen.generateRemainingKnownTypes()
	writeGeneratedCode(code, &realFS{})
}

type fileWriter interface {
	forFile(fullpath string) (io.WriteCloser, error)
}

type realFS struct{}

func (*realFS) forFile(fullpath string) (io.WriteCloser, error) {
	file, err := os.Create(fullpath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

var _ fileWriter = (*realFS)(nil)

func writeGeneratedCode(code *generatedCode, wr fileWriter) error {
	for pkg, file := range code.files {
		if len(file.impls) == 0 {
			continue
		}
		if !strings.HasPrefix(pkg, code.mod.Path) {
			log.Printf("failed to generate code for foreign package '%s'", pkg)
			log.Printf("DEBUG:\n%#v", file)
			continue
		}

		sort.Slice(file.impls, func(i, j int) bool {
			return strings.Compare(file.impls[i].name, file.impls[j].name) < 0
		})

		out := jen.NewFile(file.pkg)
		out.HeaderComment(licenseFileHeader)
		out.HeaderComment("Code generated by Sizegen. DO NOT EDIT.")

		for _, impl := range file.impls {
			if impl.flags&codeWithInterface != 0 {
				out.Add(jen.Type().Id("cachedObject").InterfaceFunc(func(i *jen.Group) {
					i.Id("CachedSize").Params(jen.Id("alloc").Id("bool")).Int64()
				}))
				break
			}
		}

		for _, impl := range file.impls {
			if impl.flags&codeWithUnsafe != 0 {
				out.Commentf("//go:nocheckptr")
			}
			out.Add(impl.code)
		}

		fullPath := path.Join(code.mod.Dir, strings.TrimPrefix(pkg, code.mod.Path), "cached_size.go")
		writer, err := wr.forFile(fullPath)
		if err != nil {
			return err
		}

		if err := out.Render(writer); err != nil {
			writer.Close()
			return fmt.Errorf("failed to save '%s': %v", fullPath, err)
		}
		if err = writer.Close(); err != nil {
			return err
		}

		log.Printf("saved %s at '%s'", pkg, fullPath)
	}
	return nil
}

func (sizegen *sizegen) generateRemainingKnownTypes() *generatedCode {
	var complete bool

	for !complete {
		complete = true
		for tt, ts := range sizegen.known {
			isComplex := !ts.pod
			notYetGenerated := !ts.generated
			if ts.local && isComplex && notYetGenerated {
				sizegen.generateKnownType(tt)
				complete = false
			}
		}
	}

	return &generatedCode{
		mod:   sizegen.mod,
		files: sizegen.codegen,
	}
}

func (sizegen *sizegen) sizeImplForStruct(name *types.TypeName, st *types.Struct) (jen.Code, codeFlag) {
	if sizegen.sizes.Sizeof(st) == 0 {
		return nil, 0
	}

	var stmt []jen.Code
	var funcFlags codeFlag
	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		fieldType := field.Type()
		fieldName := jen.Id("cached").Dot(field.Name())

		fieldStmt, flag := sizegen.sizeStmtForType(fieldName, fieldType, false)
		if fieldStmt != nil {
			if sizegen.DebugTypes {
				stmt = append(stmt, jen.Commentf("%s", field.String()))
			}
			stmt = append(stmt, fieldStmt)
		}
		funcFlags |= flag
	}

	f := jen.Func()
	f.Params(jen.Id("cached").Op("*").Id(name.Name()))
	f.Id("CachedSize").Params(jen.Id("alloc").Id("bool")).Int64()
	f.BlockFunc(func(b *jen.Group) {
		b.Add(jen.If(jen.Id("cached").Op("==").Nil()).Block(jen.Return(jen.Lit(int64(0)))))
		b.Add(jen.Id("size").Op(":=").Lit(int64(0)))
		b.Add(jen.If(jen.Id("alloc")).Block(
			jen.Id("size").Op("+=").Lit(sizegen.sizes.Sizeof(st)),
		))
		for _, s := range stmt {
			b.Add(s)
		}
		b.Add(jen.Return(jen.Id("size")))
	})
	return f, funcFlags
}

func (sizegen *sizegen) sizeStmtForMap(fieldName *jen.Statement, m *types.Map) []jen.Code {
	const bucketCnt = 8
	const sizeofHmap = int64(6 * 8)

	/*
		type bmap struct {
			// tophash generally contains the top byte of the hash value
			// for each key in this bucket. If tophash[0] < minTopHash,
			// tophash[0] is a bucket evacuation state instead.
			tophash [bucketCnt]uint8
			// Followed by bucketCnt keys and then bucketCnt elems.
			// NOTE: packing all the keys together and then all the elems together makes the
			// code a bit more complicated than alternating key/elem/key/elem/... but it allows
			// us to eliminate padding which would be needed for, e.g., map[int64]int8.
			// Followed by an overflow pointer.
		}
	*/
	sizeOfBucket := int(
		bucketCnt + // tophash
			bucketCnt*sizegen.sizes.Sizeof(m.Key()) +
			bucketCnt*sizegen.sizes.Sizeof(m.Elem()) +
			8, // overflow pointer
	)

	return []jen.Code{
		jen.Id("size").Op("+=").Lit(sizeofHmap),

		jen.Id("hmap").Op(":=").Qual("reflect", "ValueOf").Call(fieldName),

		jen.Id("numBuckets").Op(":=").Id("int").Call(
			jen.Qual("math", "Pow").Call(jen.Lit(2), jen.Id("float64").Call(
				jen.Parens(jen.Op("*").Parens(jen.Op("*").Id("uint8")).Call(
					jen.Qual("unsafe", "Pointer").Call(jen.Id("hmap").Dot("Pointer").Call().
						Op("+").Id("uintptr").Call(jen.Lit(9)))))))),

		jen.Id("numOldBuckets").Op(":=").Parens(jen.Op("*").Parens(jen.Op("*").Id("uint16")).Call(
			jen.Qual("unsafe", "Pointer").Call(
				jen.Id("hmap").Dot("Pointer").Call().Op("+").Id("uintptr").Call(jen.Lit(10))))),

		jen.Id("size").Op("+=").Id("int64").Call(jen.Id("numOldBuckets").Op("*").Lit(sizeOfBucket)),

		jen.If(jen.Id("len").Call(fieldName).Op(">").Lit(0).Op("||").Id("numBuckets").Op(">").Lit(1)).Block(
			jen.Id("size").Op("+=").Id("int64").Call(
				jen.Id("numBuckets").Op("*").Lit(sizeOfBucket))),
	}
}

func (sizegen *sizegen) sizeStmtForType(fieldName *jen.Statement, field types.Type, alloc bool) (jen.Code, codeFlag) {
	if sizegen.sizes.Sizeof(field) == 0 {
		return nil, 0
	}

	switch node := field.(type) {
	case *types.Slice:
		elemT := node.Elem()
		elemSize := sizegen.sizes.Sizeof(elemT)

		switch elemSize {
		case 0:
			return nil, 0

		case 1:
			return jen.Id("size").Op("+=").Int64().Call(jen.Cap(fieldName)), 0

		default:
			stmt, flag := sizegen.sizeStmtForType(jen.Id("elem"), elemT, false)
			return jen.BlockFunc(func(b *jen.Group) {
				b.Add(
					jen.Id("size").
						Op("+=").
						Int64().Call(jen.Cap(fieldName)).
						Op("*").
						Lit(sizegen.sizes.Sizeof(elemT)))

				if stmt != nil {
					b.Add(jen.For(jen.List(jen.Id("_"), jen.Id("elem")).Op(":=").Range().Add(fieldName)).Block(stmt))
				}
			}), flag
		}

	case *types.Map:
		keySize, keyFlag := sizegen.sizeStmtForType(jen.Id("k"), node.Key(), false)
		valSize, valFlag := sizegen.sizeStmtForType(jen.Id("v"), node.Elem(), false)

		return jen.If(fieldName.Clone().Op("!=").Nil()).BlockFunc(func(block *jen.Group) {
			for _, stmt := range sizegen.sizeStmtForMap(fieldName, node) {
				block.Add(stmt)
			}

			var forLoopVars []jen.Code
			switch {
			case keySize != nil && valSize != nil:
				forLoopVars = []jen.Code{jen.Id("k"), jen.Id("v")}
			case keySize == nil && valSize != nil:
				forLoopVars = []jen.Code{jen.Id("_"), jen.Id("v")}
			case keySize != nil && valSize == nil:
				forLoopVars = []jen.Code{jen.Id("k")}
			case keySize == nil && valSize == nil:
				return
			}

			block.Add(jen.For(jen.List(forLoopVars...).Op(":=").Range().Add(fieldName))).BlockFunc(func(b *jen.Group) {
				if keySize != nil {
					b.Add(keySize)
				}
				if valSize != nil {
					b.Add(valSize)
				}
			})
		}), codeWithUnsafe | keyFlag | valFlag

	case *types.Pointer:
		return sizegen.sizeStmtForType(fieldName, node.Elem(), true)

	case *types.Named:
		ts := sizegen.getKnownType(node)
		if ts.pod || !ts.local {
			if alloc {
				if !ts.local {
					log.Printf("WARNING: size of external type %s cannot be fully calculated", node)
				}
				return jen.If(fieldName.Clone().Op("!=").Nil()).Block(
					jen.Id("size").Op("+=").Lit(sizegen.sizes.Sizeof(node.Underlying())),
				), 0
			}
			return nil, 0
		}
		return sizegen.sizeStmtForType(fieldName, node.Underlying(), alloc)

	case *types.Interface:
		if node.Empty() {
			return nil, 0
		}
		return jen.If(
			jen.List(
				jen.Id("cc"), jen.Id("ok")).
				Op(":=").
				Add(fieldName.Clone().Assert(jen.Id("cachedObject"))),
			jen.Id("ok"),
		).Block(
			jen.Id("size").
				Op("+=").
				Id("cc").
				Dot("CachedSize").
				Call(jen.True()),
		), codeWithInterface

	case *types.Struct:
		return jen.Id("size").Op("+=").Add(fieldName.Clone().Dot("CachedSize").Call(jen.Lit(alloc))), 0

	case *types.Basic:
		if !alloc {
			if node.Info()&types.IsString != 0 {
				return jen.Id("size").Op("+=").Int64().Call(jen.Len(fieldName)), 0
			}
			return nil, 0
		}
		return jen.Id("size").Op("+=").Lit(sizegen.sizes.Sizeof(node)), 0
	default:
		log.Printf("unhandled type: %T", node)
		return nil, 0
	}
}

type typePaths []string

func (t *typePaths) String() string {
	return fmt.Sprintf("%v", *t)
}

func (t *typePaths) Set(path string) error {
	*t = append(*t, path)
	return nil
}

func main() {
	var patterns typePaths
	var generate typePaths
	flag.Var(&patterns, "in", "Go packages to load the generator")
	flag.Var(&generate, "gen", "Typename of the Go struct to generate size info for")
	flag.Parse()

	loaded, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedTypesInfo | packages.NeedDeps | packages.NeedImports | packages.NeedModule,
		Logf: log.Printf,
	}, patterns...)

	if err != nil {
		log.Fatal(err)
	}

	sizegen, err := generateCode(loaded, generate)
	if err != nil {
		log.Fatal(err)
	}

	sizegen.finalize()
}

func generateCode(loaded []*packages.Package, generate typePaths) (*sizegen, error) {
	sizegen := newSizegen(loaded[0].Module, loaded[0].TypesSizes)

	scopes := make(map[string]*types.Scope)
	for _, pkg := range loaded {
		scopes[pkg.PkgPath] = pkg.Types.Scope()
	}

	for _, gen := range generate {
		pos := strings.LastIndexByte(gen, '.')
		if pos < 0 {
			return nil, fmt.Errorf("unexpected input type: %s", gen)
		}

		pkgname := gen[:pos]
		typename := gen[pos+1:]

		scope := scopes[pkgname]
		if scope == nil {
			return nil, fmt.Errorf("no scope found for type '%s'", gen)
		}

		tt := scope.Lookup(typename)
		if tt == nil {
			return nil, fmt.Errorf("no type called '%s' found in '%s'", typename, pkgname)
		}

		sizegen.generateKnownType(tt.Type().(*types.Named))
	}

	return sizegen, nil
}
