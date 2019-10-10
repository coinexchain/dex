package codon

import (
	"crypto/sha256"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"unicode"
)

const (
	VersionOfAPI = 1
)

func ShowInfoForVar(leafTypes map[string]string, v interface{}) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	fmt.Printf("======= %v '%s' '%s' == \n", t, t.PkgPath(), t.Name())
	showInfo(leafTypes, "", t)
}

func structHasPrivateField(t reflect.Type) bool {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		var isPrivate bool
		for _, r := range field.Name {
			isPrivate = unicode.IsLower(r)
			break
		}
		if isPrivate {
			return true
		}
	}
	return false
}

func showInfo(leafTypes map[string]string, indent string, t reflect.Type) {
	ending := ""
	indentP := indent + "    "
	switch t.Kind() {
	case reflect.Bool:
		fmt.Printf("bool")
	case reflect.Int:
		fmt.Printf("int")
	case reflect.Int8:
		fmt.Printf("int8")
	case reflect.Int16:
		fmt.Printf("int16")
	case reflect.Int32:
		fmt.Printf("int32")
	case reflect.Int64:
		fmt.Printf("int64")
	case reflect.Uint:
		fmt.Printf("uint")
	case reflect.Uint8:
		fmt.Printf("uint8")
	case reflect.Uint16:
		fmt.Printf("uint16")
	case reflect.Uint32:
		fmt.Printf("uint32")
	case reflect.Uint64:
		fmt.Printf("uint64")
	case reflect.Uintptr:
		fmt.Printf("Uintptr!")
	case reflect.Complex64:
		fmt.Printf("complex64!")
	case reflect.Complex128:
		fmt.Printf("complex128!")
	case reflect.Float32:
		fmt.Printf("float32")
	case reflect.Float64:
		fmt.Printf("float64")
	case reflect.Chan:
		fmt.Printf("chan!")
	case reflect.Func:
		fmt.Printf("func!")
	case reflect.Interface:
		fmt.Printf("interface (%s %s)!", t.PkgPath(), t.Name())
	case reflect.Map:
		fmt.Printf("map!")
	case reflect.Ptr:
		path := t.Elem().PkgPath() + "." + t.Elem().Name()
		if _, ok := leafTypes[path]; ok {
			fmt.Printf("pointer ('%s' '%s')\n", t.Elem().PkgPath(), t.Elem().Name())
		} else {
			fmt.Printf("pointer ('%s' '%s') {\n", t.Elem().PkgPath(), t.Elem().Name())
			fmt.Printf("%s", indentP)
			showInfo(leafTypes, indentP, t.Elem())
			ending = indent + "} // pointer"
		}
	case reflect.Array:
		fmt.Printf("array {\n")
		fmt.Printf("%s", indentP)
		showInfo(leafTypes, indentP, t.Elem())
		ending = indent + "} //array"
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			fmt.Printf("ByteSlice")
		} else {
			fmt.Printf("slice {\n")
			fmt.Printf("%s", indentP)
			showInfo(leafTypes, indentP, t.Elem())
			ending = indent + "} //slice"
		}
	case reflect.String:
		fmt.Printf("string")
	case reflect.Struct:
		path := t.PkgPath() + "." + t.Name()
		if _, ok := leafTypes[path]; ok {
			fmt.Printf("struct ('%s' '%s')\n", t.PkgPath(), t.Name())
		} else {
			if structHasPrivateField(t) {
				fmt.Printf("struct_with_private {\n")
			} else {
				fmt.Printf("struct {\n")
			}
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				fmt.Printf("%s%s : ('%s' '%s') ", indentP, field.Name, field.Type.PkgPath(), field.Type.Name())
				path = field.Type.PkgPath() + "." + field.Type.Name()
				if _, ok := leafTypes[path]; ok {
					fmt.Printf("\n")
				} else {
					showInfo(leafTypes, indentP, field.Type)
				}
			}
			ending = indent + "} //struct"
		}
	default:
		fmt.Printf("Unknown Kind! %s", t.Kind())
	}

	fmt.Printf("%s\n", ending)
}

type MagicBytes [4]byte

func calcMagicBytes(lines []string) [4]byte {
	var res [4]byte
	h := sha256.New()
	for _, line := range lines {
		h.Write([]byte(line))
	}
	bz := h.Sum(nil)
	for i := 0; i < 4; i++ {
		res[i] = bz[i]
	}
	return res
}

type AliasAndValue struct {
	Alias string
	Value interface{}
}

func writeLines(w io.Writer, lines []string) {
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		w.Write([]byte(line))
		w.Write([]byte("\n"))
	}
}

func GenerateCodecFile(w io.Writer, leafTypes, ignoreImpl map[string]string, aliasAndValueList []AliasAndValue,
	extraLogics string, extraImports []string) {

	w.Write([]byte("//nolint\npackage codec\nimport (\n"))
	for _, p := range extraImports {
		w.Write([]byte(p + "\n"))
	}
	w.Write([]byte("\"io\"\n\"encoding/binary\"\n\"math\"\n\"errors\"\n)\n"))
	w.Write([]byte(headerLogics))
	w.Write([]byte(extraLogics))
	ctx := newPrepareCtx(leafTypes, ignoreImpl)
	for _, entry := range aliasAndValueList {
		ctx.register(entry.Alias, entry.Value)
	}
	ctx.analyzeIfc()
	for _, entry := range aliasAndValueList {
		t := derefPtr(entry.Value)
		if t.Kind() != reflect.Interface {
			w.Write([]byte("// Non-Interface\n"))
			lines := ctx.prepareStructFunc(entry.Alias, t)
			writeLines(w, lines)
		}
	}
	for _, entry := range aliasAndValueList {
		t := derefPtr(entry.Value)
		if t.Kind() == reflect.Interface {
			w.Write([]byte("// Interface\n"))
			lines := ctx.prepareIfcFunc(entry.Alias, t)
			writeLines(w, lines)
		}
	}
	lines := ctx.prepareMagicBytesFunc()
	writeLines(w, lines)

	aliases := make([]string, 0, len(ctx.structPath2Alias))
	for _, alias := range ctx.structPath2Alias {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	lines = prepareIfcEncodeFunc("EncodeAny", aliases)
	writeLines(w, lines)
	lines = prepareBareEncodeAnyFunc(aliases)
	writeLines(w, lines)
	lines = ctx.prepareDecodeAnyFunc()
	writeLines(w, lines)
	lines = prepareBareDecodeAnyFunc(aliases)
	writeLines(w, lines)
	lines = prepareIfcRandFunc("RandAny", "interface{}", aliases, nil)
	writeLines(w, lines)
	lines = ctx.prepareSupportListFunc()
	writeLines(w, lines)
}

type prepareCtx struct {
	structPath2Alias map[string]string
	ifcPath2Alias    map[string]string
	structPath2Type  map[string]reflect.Type
	ifcPath2Type     map[string]reflect.Type

	ifcPath2StructPaths map[string][]string

	structAlias2MagicBytes map[string]MagicBytes
	magicBytes2Alias       map[MagicBytes]string

	leafTypes  map[string]string
	ignoreImpl map[string]string
}

func newPrepareCtx(leafTypes, ignoreImpl map[string]string) *prepareCtx {
	return &prepareCtx{
		structPath2Alias: make(map[string]string),
		ifcPath2Alias:    make(map[string]string),
		structPath2Type:  make(map[string]reflect.Type),
		ifcPath2Type:     make(map[string]reflect.Type),

		ifcPath2StructPaths:    make(map[string][]string),
		structAlias2MagicBytes: make(map[string]MagicBytes),
		magicBytes2Alias:       make(map[MagicBytes]string),
		leafTypes:              leafTypes,
		ignoreImpl:             ignoreImpl,
	}
}

func prepareIfcEncodeFunc(funcName string, aliases []string) []string {
	lines := make([]string, 0, 1000)
	lines = append(lines, "func "+funcName+"(w io.Writer, x interface{}) error {")
	lines = append(lines, "switch v := x.(type) {")
	for _, alias := range aliases {
		lines = append(lines, fmt.Sprintf("case %s:", alias))
		lines = append(lines, fmt.Sprintf("w.Write(getMagicBytes(\"%s\"))", alias))
		lines = append(lines, fmt.Sprintf("return Encode%s(w, v)", alias))

		lines = append(lines, fmt.Sprintf("case *%s:", alias))
		lines = append(lines, fmt.Sprintf("w.Write(getMagicBytes(\"%s\"))", alias))
		lines = append(lines, fmt.Sprintf("return Encode%s(w, *v)", alias))
	}
	lines = append(lines, "default:")
	lines = append(lines, "panic(\"Unknown Type.\")")
	lines = append(lines, "} // end of switch")
	lines = append(lines, "} // end of func")
	return lines
}

func prepareIfcRandFunc(funcName, ifc string, aliases []string, ignoreImpl map[string]string) []string {
	lines := make([]string, 0, 1000)
	lines = append(lines, "func "+funcName+"(r RandSrc) "+ifc+" {")
	newAliases := make([]string, 0, len(aliases))
	for _, alias := range aliases {
		if ignoreImpl == nil || ignoreImpl[alias] != ifc {
			newAliases = append(newAliases, alias)
		}
	}
	lines = append(lines, fmt.Sprintf("switch r.GetUint() %% %d {", len(newAliases)))
	for i, alias := range newAliases {
		lines = append(lines, fmt.Sprintf("case %d:", i))
		lines = append(lines, fmt.Sprintf("return Rand%s(r)", alias))
	}
	lines = append(lines, "default:")
	lines = append(lines, "panic(\"Unknown Type.\")")
	lines = append(lines, "} // end of switch")
	lines = append(lines, "} // end of func")
	return lines
}

func (ctx *prepareCtx) prepareDecodeAnyFunc() []string {
	res, _ := prepareIfcDecodeFunc("DecodeAny", "interface{}", ctx.structAlias2MagicBytes)
	return res
}

func prepareIfcDecodeFunc(funcName, decType string, alias2bytes map[string]MagicBytes) ([]string, []string) {
	aliases := make([]string, 0, len(alias2bytes))
	for alias := range alias2bytes {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	lines := make([]string, 0, 1000)
	lines = append(lines, "func "+funcName+"(bz []byte) ("+decType+", int, error) {")
	lines = append(lines, "var v "+decType)
	lines = append(lines, "var magicBytes [4]byte")
	lines = append(lines, "var n int")
	lines = append(lines, "for i:=0; i<4; i++ {magicBytes[i] = bz[i]}")
	lines = append(lines, "switch magicBytes {")
	for _, alias := range aliases {
		magicBytes := alias2bytes[alias]
		lines = append(lines, fmt.Sprintf("case [4]byte{%d,%d,%d,%d}:",
			magicBytes[0], magicBytes[1], magicBytes[2], magicBytes[3]))
		lines = append(lines, fmt.Sprintf("v, n, err := Decode%s(bz[4:])", alias))
		lines = append(lines, fmt.Sprintf("return v, n+4, err"))
	}
	lines = append(lines, "default:")
	lines = append(lines, "panic(\"Unknown type\")")
	lines = append(lines, "} // end of switch")
	lines = append(lines, "return v, n, nil")
	lines = append(lines, "} // end of "+funcName)
	return lines, aliases
}

func prepareBareEncodeAnyFunc(aliases []string) []string {
	lines := make([]string, 0, 1000)
	lines = append(lines, "func BareEncodeAny(w io.Writer, x interface{}) error {")
	lines = append(lines, "switch v := x.(type) {")
	for _, alias := range aliases {
		lines = append(lines, fmt.Sprintf("case %s:", alias))
		lines = append(lines, fmt.Sprintf("return Encode%s(w, v)", alias))

		lines = append(lines, fmt.Sprintf("case *%s:", alias))
		lines = append(lines, fmt.Sprintf("return Encode%s(w, *v)", alias))
	}
	lines = append(lines, "default:")
	lines = append(lines, "panic(\"Unknown Type.\")")
	lines = append(lines, "} // end of switch")
	lines = append(lines, "} // end of func")
	return lines
}

func prepareBareDecodeAnyFunc(aliases []string) []string {
	lines := make([]string, 0, 1000)
	lines = append(lines, "func BareDecodeAny(bz []byte, x interface{}) (n int, err error) {")
	lines = append(lines, "switch v := x.(type) {")
	for _, alias := range aliases {
		lines = append(lines, fmt.Sprintf("case *%s:", alias))
		lines = append(lines, fmt.Sprintf("*v, n, err = Decode%s(bz)", alias))
	}
	lines = append(lines, "default:")
	lines = append(lines, "panic(\"Unknown type\")")
	lines = append(lines, "} // end of switch")
	lines = append(lines, "return")
	lines = append(lines, "} // end of DecodeVar")
	return lines
}

func (ctx *prepareCtx) prepareMagicBytesFunc() []string {
	lines := make([]string, 0, 1000)
	lines = append(lines, "func getMagicBytes(name string) []byte {")
	lines = append(lines, "switch name {")
	aliases := make([]string, 0, len(ctx.structAlias2MagicBytes))
	for alias := range ctx.structAlias2MagicBytes {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	for _, alias := range aliases {
		magicBytes := ctx.structAlias2MagicBytes[alias]
		lines = append(lines, fmt.Sprintf("case \"%s\":", alias))
		lines = append(lines, fmt.Sprintf("return []byte{%d,%d,%d,%d}",
			magicBytes[0], magicBytes[1], magicBytes[2], magicBytes[3]))
	}
	lines = append(lines, "} // end of switch")
	lines = append(lines, "panic(\"Should not reach here\")")
	lines = append(lines, "return []byte{}")
	lines = append(lines, "} // end of getMagicBytes")
	return lines
}

func (ctx *prepareCtx) prepareSupportListFunc() []string {
	length := len(ctx.structPath2Alias) + len(ctx.ifcPath2Alias) + 10
	paths := make([]string, 0, length)
	for path := range ctx.structPath2Alias {
		paths = append(paths, path)
	}
	for path := range ctx.ifcPath2Alias {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	lines := make([]string, 0, length)
	lines = append(lines, "func GetSupportList() []string {")
	lines = append(lines, "return []string {")
	for _, path := range paths {
		lines = append(lines, fmt.Sprintf("\"%s\",", path))
	}
	lines = append(lines, "}")
	lines = append(lines, "} // end of GetSupportList")
	return lines
}

func (ctx *prepareCtx) analyzeIfc() {
	for ifcPath, ifcType := range ctx.ifcPath2Type {
		for structPath, structType := range ctx.structPath2Type {
			if structType.Implements(ifcType) {
				if _, ok := ctx.ifcPath2StructPaths[ifcPath]; ok {
					ctx.ifcPath2StructPaths[ifcPath] = append(ctx.ifcPath2StructPaths[ifcPath], structPath)
				} else {
					ctx.ifcPath2StructPaths[ifcPath] = []string{structPath}
				}
			}
		}
	}
}

func derefPtr(v interface{}) reflect.Type {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func (ctx *prepareCtx) register(alias string, v interface{}) {
	t := derefPtr(v)
	path := t.PkgPath() + "." + t.Name()
	if len(t.PkgPath()) == 0 || len(t.Name()) == 0 {
		panic("Invalid Path:" + path)
	}
	if t.Kind() == reflect.Interface {
		ctx.ifcPath2Alias[path] = alias
		ctx.ifcPath2Type[path] = t
	} else {
		ctx.structPath2Alias[path] = alias
		ctx.structPath2Type[path] = t
	}
}

func (ctx *prepareCtx) prepareIfcFunc(ifc string, t reflect.Type) []string {
	ifcPath := t.PkgPath() + "." + t.Name()
	structPaths, ok := ctx.ifcPath2StructPaths[ifcPath]
	if !ok {
		panic("Cannot find implementations for " + ifc)
	}
	alias2bytes := make(map[string]MagicBytes, len(structPaths))
	for _, structPath := range structPaths {
		alias, ok := ctx.structPath2Alias[structPath]
		if !ok {
			panic("Cannot find alias")
		}
		magicBytes, ok := ctx.structAlias2MagicBytes[alias]
		if !ok {
			panic("Cannot find magicbytes")
		}
		alias2bytes[alias] = magicBytes
	}
	decLines, aliases := prepareIfcDecodeFunc("Decode"+ifc, ifc, alias2bytes)
	encLines := prepareIfcEncodeFunc("Encode"+ifc, aliases)
	randLines := prepareIfcRandFunc("Rand"+ifc, ifc, aliases, ctx.ignoreImpl)
	return append(append(encLines, decLines...), randLines...)
}

func (ctx *prepareCtx) prepareStructFunc(alias string, t reflect.Type) []string {
	lines := make([]string, 0, 1000)

	apiLine := fmt.Sprintf("// codon version: %d", VersionOfAPI)
	line := fmt.Sprintf("func Encode%s(w io.Writer, v %s) error {", alias, alias)
	lines = append(lines, line)
	lines = append(lines, apiLine)
	lines = append(lines, "var err error")
	if t.Kind() == reflect.Struct {
		ctx.genStructEncLines(t, &lines, "v", 0)
	} else {
		ctx.genFieldEncLines(t, &lines, "v", 0)
	}
	lines = append(lines, "return nil")
	lines = append(lines, "} //End of Encode"+alias+"\n")

	line = fmt.Sprintf("func Decode%s(bz []byte) (%s, int, error) {", alias, alias)
	lines = append(lines, line)
	lines = append(lines, apiLine)
	lines = append(lines, "var err error")
	lengthLinePosition := len(lines)
	lines = append(lines, "") // length placeholder
	lines = append(lines, "var v "+alias)
	lines = append(lines, "var n int")
	lines = append(lines, "var total int")
	needLength := false
	if t.Kind() == reflect.Struct {
		nl := ctx.genStructDecLines(t, &lines, "v", 0)
		needLength = needLength || nl
	} else {
		nl := ctx.genFieldDecLines(t, &lines, "v", 0)
		needLength = needLength || nl
	}
	if needLength {
		lines[lengthLinePosition] = "var length int"
	}
	lines = append(lines, "return v, total, nil")
	lines = append(lines, "} //End of Decode"+alias+"\n")

	line = fmt.Sprintf("func Rand%s(r RandSrc) %s {", alias, alias)
	lines = append(lines, line)
	lines = append(lines, apiLine)
	lengthLinePosition = len(lines)
	lines = append(lines, "") // length placeholder
	lines = append(lines, "var v "+alias)
	needLength = false
	if t.Kind() == reflect.Struct {
		nl := ctx.genStructRandLines(t, &lines, "v", 0)
		needLength = needLength || nl
	} else {
		nl := ctx.genFieldRandLines(t, &lines, "v", 0)
		needLength = needLength || nl
	}
	if needLength {
		lines[lengthLinePosition] = "var length int"
	}
	lines = append(lines, "return v")
	lines = append(lines, "} //End of Rand"+alias+"\n")

	magicBytes := calcMagicBytes(lines)
	if otherAlias, ok := ctx.magicBytes2Alias[magicBytes]; ok {
		panic("Magic Bytes Conflicts: " + otherAlias + " vs " + alias)
	}
	ctx.structAlias2MagicBytes[alias] = magicBytes
	ctx.magicBytes2Alias[magicBytes] = alias
	return lines
}

func (ctx *prepareCtx) genFieldEncLines(t reflect.Type, lines *[]string, fieldName string, iterLevel int) {
	ending := "\nif err != nil {return err}"
	isPtr := false
	if t.Kind() == reflect.Ptr {
		isPtr = true
		elemT := t.Elem()
		if elemT.Kind() == reflect.Struct {
			t = elemT
		} else {
			panic(fmt.Sprintf("Pointer to %s is not supported", elemT.Kind()))
		}
	}
	var line string
	switch t.Kind() {
	case reflect.Chan:
		panic("Channel is not supported")
	case reflect.Func:
		panic("Func is not supported")
	case reflect.Uintptr:
		panic("Uintptr is not supported")
	case reflect.Complex64:
		panic("Complex64 is not supported")
	case reflect.Complex128:
		panic("Complex128 is not supported")
	case reflect.Map:
		panic("Map is not supported")

	case reflect.Bool:
		if len(t.PkgPath()) == 0 {
			line = fmt.Sprintf("err = codonEncodeBool(w, %s)%s", fieldName, ending)
		} else {
			line = fmt.Sprintf("err = codonEncodeBool(w, bool(%s))%s", fieldName, ending)
		}
	case reflect.Int:
		line = fmt.Sprintf("err = codonEncodeVarint(w, int64(%s))%s", fieldName, ending)
	case reflect.Int8:
		if len(t.PkgPath()) == 0 {
			line = fmt.Sprintf("err = codonEncodeInt8(w, %s)%s", fieldName, ending)
		} else {
			line = fmt.Sprintf("err = codonEncodeInt8(w, int8(%s))%s", fieldName, ending)
		}
	case reflect.Int16:
		if len(t.PkgPath()) == 0 {
			line = fmt.Sprintf("err = codonEncodeInt16(w, %s)%s", fieldName, ending)
		} else {
			line = fmt.Sprintf("err = codonEncodeInt16(w, int16(%s))%s", fieldName, ending)
		}
	case reflect.Int32:
		line = fmt.Sprintf("err = codonEncodeVarint(w, int64(%s))%s", fieldName, ending)
	case reflect.Int64:
		line = fmt.Sprintf("err = codonEncodeVarint(w, int64(%s))%s", fieldName, ending)
	case reflect.Uint:
		line = fmt.Sprintf("err = codonEncodeUvarint(w, uint64(%s))%s", fieldName, ending)
	case reflect.Uint8:
		if len(t.PkgPath()) == 0 {
			line = fmt.Sprintf("err = codonEncodeUint8(w, %s)%s", fieldName, ending)
		} else {
			line = fmt.Sprintf("err = codonEncodeUint8(w, uint8(%s))%s", fieldName, ending)
		}
	case reflect.Uint16:
		if len(t.PkgPath()) == 0 {
			line = fmt.Sprintf("err = codonEncodeUint16(w, %s)%s", fieldName, ending)
		} else {
			line = fmt.Sprintf("err = codonEncodeUint16(w, uint16(%s))%s", fieldName, ending)
		}
	case reflect.Uint32:
		line = fmt.Sprintf("err = codonEncodeUvarint(w, uint64(%s))%s", fieldName, ending)
	case reflect.Uint64:
		line = fmt.Sprintf("err = codonEncodeUvarint(w, uint64(%s))%s", fieldName, ending)
	case reflect.Float32:
		if len(t.PkgPath()) == 0 {
			line = fmt.Sprintf("err = codonEncodeFloat32(w, %s)%s", fieldName, ending)
		} else {
			line = fmt.Sprintf("err = codonEncodeFloat32(w, float32(%s))%s", fieldName, ending)
		}
	case reflect.Float64:
		if len(t.PkgPath()) == 0 {
			line = fmt.Sprintf("err = codonEncodeFloat64(w, %s)%s", fieldName, ending)
		} else {
			line = fmt.Sprintf("err = codonEncodeFloat64(w, float64(%s))%s", fieldName, ending)
		}
	case reflect.String:
		if len(t.PkgPath()) == 0 {
			line = fmt.Sprintf("err = codonEncodeString(w, %s)%s", fieldName, ending)
		} else {
			line = fmt.Sprintf("err = codonEncodeString(w, string(%s))%s", fieldName, ending)
		}
	case reflect.Array, reflect.Slice:
		elemT := t.Elem()
		if elemT.Kind() == reflect.Uint8 {
			line = fmt.Sprintf("err = codonEncodeByteSlice(w, %s[:])%s", fieldName, ending)
		} else {
			line = fmt.Sprintf("err = codonEncodeVarint(w, int64(len(%s)))%s", fieldName, ending)
			*lines = append(*lines, line)
			iterVar := fmt.Sprintf("_%d", iterLevel)
			line = fmt.Sprintf("for %s:=0; %s<len(%s); %s++ {",
				iterVar, iterVar, fieldName, iterVar)
			*lines = append(*lines, line)
			varName := fieldName + "[" + iterVar + "]"
			ctx.genFieldEncLines(elemT, lines, varName, iterLevel+1)
			line = "}"
		}
	case reflect.Interface:
		typePath := t.PkgPath() + "." + t.Name()
		alias, ok := ctx.ifcPath2Alias[typePath]
		if !ok {
			panic("Cannot find alias for:" + typePath)
		}
		line = fmt.Sprintf("err = Encode%s(w, %s)%s // interface_encode", alias, fieldName, ending)
	case reflect.Ptr:
		panic("Should not reach here")
	case reflect.Struct:
		if _, ok := ctx.leafTypes[t.PkgPath()+"."+t.Name()]; ok {
			if isPtr {
				fieldName = "*(" + fieldName + ")"
			}
			line = fmt.Sprintf("err = Encode%s(w, %s)%s", t.Name(), fieldName, ending)
		} else {
			ctx.genStructEncLines(t, lines, fieldName, iterLevel)
			line = "// end of " + fieldName
		}
	default:
		panic(fmt.Sprintf("Unknown Kind %s", t.Kind()))
	}
	*lines = append(*lines, line)
}

func (ctx *prepareCtx) genStructEncLines(t reflect.Type, lines *[]string, varName string, iterLevel int) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		ctx.genFieldEncLines(field.Type, lines, varName+"."+field.Name, iterLevel)
	}
}

//=========================

func (ctx *prepareCtx) getTypeName(elemT reflect.Type) string {
	if elemT.Kind() == reflect.Ptr {
		panic("slice/array of pointers are not support")
	}
	if len(elemT.PkgPath()) == 0 {
		return elemT.Name() //basic type
	}
	typePath := elemT.PkgPath() + "." + elemT.Name()
	alias, ok := ctx.structPath2Alias[typePath]
	if !ok {
		alias, ok = ctx.ifcPath2Alias[typePath]
	}
	if !ok {
		panic(typePath + " is not registered")
	}
	return alias
}

func (ctx *prepareCtx) buildDecLine(typeName, fieldName, ending string, t reflect.Type) string {
	if len(t.PkgPath()) == 0 {
		return fmt.Sprintf("%s = %s(codonDecode%s(bz, &n, &err))%s", fieldName, strings.ToLower(typeName), typeName, ending)
	}
	alias := ctx.getTypeName(t)
	return fmt.Sprintf("%s = %s(codonDecode%s(bz, &n, &err))%s", fieldName, alias, typeName, ending)
}

func (ctx *prepareCtx) genFieldDecLines(t reflect.Type, lines *[]string, fieldName string, iterLevel int) bool {
	ending := "\nif err != nil {return v, total, err}\nbz = bz[n:]\ntotal+=n"
	needLength := false
	isPtr := false
	if t.Kind() == reflect.Ptr {
		isPtr = true
		elemT := t.Elem()
		if elemT.Kind() == reflect.Struct {
			t = elemT
		} else {
			panic(fmt.Sprintf("Pointer to %s is not supported", elemT.Kind()))
		}
	}
	var line string
	switch t.Kind() {
	case reflect.Chan:
		panic("Channel is not supported")
	case reflect.Func:
		panic("Func is not supported")
	case reflect.Uintptr:
		panic("Uintptr is not supported")
	case reflect.Complex64:
		panic("Complex64 is not supported")
	case reflect.Complex128:
		panic("Complex128 is not supported")
	case reflect.Map:
		panic("Map is not supported")
	case reflect.Bool:
		line = ctx.buildDecLine("Bool", fieldName, ending, t)
	case reflect.Int:
		line = ctx.buildDecLine("Int", fieldName, ending, t)
	case reflect.Int8:
		line = ctx.buildDecLine("Int8", fieldName, ending, t)
	case reflect.Int16:
		line = ctx.buildDecLine("Int16", fieldName, ending, t)
	case reflect.Int32:
		line = ctx.buildDecLine("Int32", fieldName, ending, t)
	case reflect.Int64:
		line = ctx.buildDecLine("Int64", fieldName, ending, t)
	case reflect.Uint:
		line = ctx.buildDecLine("Uint", fieldName, ending, t)
	case reflect.Uint8:
		line = ctx.buildDecLine("Uint8", fieldName, ending, t)
	case reflect.Uint16:
		line = ctx.buildDecLine("Uint16", fieldName, ending, t)
	case reflect.Uint32:
		line = ctx.buildDecLine("Uint32", fieldName, ending, t)
	case reflect.Uint64:
		line = ctx.buildDecLine("Uint64", fieldName, ending, t)
	case reflect.Float32:
		line = ctx.buildDecLine("Float32", fieldName, ending, t)
	case reflect.Float64:
		line = ctx.buildDecLine("Float64", fieldName, ending, t)
	case reflect.String:
		line = ctx.buildDecLine("String", fieldName, ending, t)
	case reflect.Array, reflect.Slice:
		line = fmt.Sprintf("length = codonDecodeInt(bz, &n, &err)%s", ending)
		needLength = true
		*lines = append(*lines, line)
		typeName := ctx.getTypeName(t.Elem())
		elemT := t.Elem()
		if t.Kind() == reflect.Slice && elemT.Kind() != reflect.Uint8 {
			makeSlice := fmt.Sprintf("%s = make([]%s, length)", fieldName, typeName)
			*lines = append(*lines, makeSlice)
		}
		if elemT.Kind() == reflect.Uint8 && t.Kind() == reflect.Slice {
			line = fmt.Sprintf("%s, n, err = codonGetByteSlice(bz, length)%s", fieldName, ending)
		} else {
			iterVar := fmt.Sprintf("_%d", iterLevel)
			initVar := fmt.Sprintf("%s, length_%d := 0, length", iterVar, iterLevel)
			line = fmt.Sprintf("for %s; %s<length_%d; %s++ { //%s of %s",
				initVar, iterVar, iterLevel, iterVar, t.Kind(), t.Elem().Kind())
			*lines = append(*lines, line)
			if t.Elem().Kind() == reflect.Interface || t.Elem().Kind() == reflect.Struct {
				line = fmt.Sprintf("%s[%s], n, err = Decode%s(bz)%s", fieldName, iterVar, typeName, ending)
				*lines = append(*lines, line)
			} else {
				varName := fieldName + "[" + iterVar + "]"
				nl := ctx.genFieldDecLines(elemT, lines, varName, iterLevel+1)
				needLength = needLength || nl
			}
			line = "}"
		}
	case reflect.Interface:
		typePath := t.PkgPath() + "." + t.Name()
		alias, ok := ctx.ifcPath2Alias[typePath]
		if !ok {
			panic("Cannot find alias for:" + typePath)
		}
		line = fmt.Sprintf("%s, n, err = Decode%s(bz)%s // interface_decode", fieldName, alias, ending)
	case reflect.Ptr:
		panic("Should not reach here")
	case reflect.Struct:
		if _, ok := ctx.leafTypes[t.PkgPath()+"."+t.Name()]; ok {
			if isPtr {
				*lines = append(*lines, ctx.initPtrMember(fieldName, t))
				line = fmt.Sprintf("*(%s), n, err = Decode%s(bz)%s", fieldName, t.Name(), ending)
			} else {
				line = fmt.Sprintf("%s, n, err = Decode%s(bz)%s", fieldName, t.Name(), ending)
			}
		} else {
			if isPtr {
				*lines = append(*lines, ctx.initPtrMember(fieldName, t))
			}
			nl := ctx.genStructDecLines(t, lines, fieldName, iterLevel)
			needLength = needLength || nl
			line = "// end of " + fieldName
		}
	default:
		panic(fmt.Sprintf("Unknown Kind %s", t.Kind()))
	}
	*lines = append(*lines, line)
	return needLength
}

func (ctx *prepareCtx) initPtrMember(fieldName string, t reflect.Type) string {
	typePath := t.PkgPath() + "." + t.Name()
	alias, ok := ctx.structPath2Alias[typePath]
	if !ok {
		alias, ok = ctx.leafTypes[typePath]
	}
	if !ok {
		panic("Cannot find alias for:" + typePath)
	}
	return fmt.Sprintf("%s = &%s{}", fieldName, alias)
}

func (ctx *prepareCtx) genStructDecLines(t reflect.Type, lines *[]string, varName string, iterLevel int) bool {
	needLength := false
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		nl := ctx.genFieldDecLines(field.Type, lines, varName+"."+field.Name, iterLevel)
		needLength = needLength || nl
	}
	return needLength
}

//======================

func (ctx *prepareCtx) buildRandLine(typeName, fieldName string, t reflect.Type) string {
	if len(t.PkgPath()) == 0 {
		return fmt.Sprintf("%s = r.Get%s()", fieldName, typeName)
	}
	alias := ctx.getTypeName(t)
	return fmt.Sprintf("%s = %s(r.Get%s())", fieldName, alias, typeName)
}

func (ctx *prepareCtx) genFieldRandLines(t reflect.Type, lines *[]string, fieldName string, iterLevel int) bool {
	needLength := false
	isPtr := false
	if t.Kind() == reflect.Ptr {
		isPtr = true
		elemT := t.Elem()
		if elemT.Kind() == reflect.Struct {
			t = elemT
		} else {
			panic(fmt.Sprintf("Pointer to %s is not supported", elemT.Kind()))
		}
	}
	var line string
	switch t.Kind() {
	case reflect.Chan:
		panic("Channel is not supported")
	case reflect.Func:
		panic("Func is not supported")
	case reflect.Uintptr:
		panic("Uintptr is not supported")
	case reflect.Complex64:
		panic("Complex64 is not supported")
	case reflect.Complex128:
		panic("Complex128 is not supported")
	case reflect.Map:
		panic("Map is not supported")
	case reflect.Bool:
		line = ctx.buildRandLine("Bool", fieldName, t)
	case reflect.Int:
		line = ctx.buildRandLine("Int", fieldName, t)
	case reflect.Int8:
		line = ctx.buildRandLine("Int8", fieldName, t)
	case reflect.Int16:
		line = ctx.buildRandLine("Int16", fieldName, t)
	case reflect.Int32:
		line = ctx.buildRandLine("Int32", fieldName, t)
	case reflect.Int64:
		line = ctx.buildRandLine("Int64", fieldName, t)
	case reflect.Uint:
		line = ctx.buildRandLine("Uint", fieldName, t)
	case reflect.Uint8:
		line = ctx.buildRandLine("Uint8", fieldName, t)
	case reflect.Uint16:
		line = ctx.buildRandLine("Uint16", fieldName, t)
	case reflect.Uint32:
		line = ctx.buildRandLine("Uint32", fieldName, t)
	case reflect.Uint64:
		line = ctx.buildRandLine("Uint64", fieldName, t)
	case reflect.Float32:
		line = ctx.buildRandLine("Float32", fieldName, t)
	case reflect.Float64:
		line = ctx.buildRandLine("Float64", fieldName, t)
	case reflect.String:
		line = fmt.Sprintf("%s = r.GetString(1+int(r.GetUint()%%(MaxStringLength-1)))", fieldName)
	case reflect.Array, reflect.Slice:
		line = "length = 1+int(r.GetUint()%(MaxSliceLength-1))"
		if t.Kind() == reflect.Array {
			line = fmt.Sprintf("length = %d", t.Len())
		}
		needLength = true
		*lines = append(*lines, line)
		typeName := ctx.getTypeName(t.Elem())
		elemT := t.Elem()
		if t.Kind() == reflect.Slice && elemT.Kind() != reflect.Uint8 {
			makeSlice := fmt.Sprintf("%s = make([]%s, length)", fieldName, typeName)
			*lines = append(*lines, makeSlice)
		}
		if elemT.Kind() == reflect.Uint8 && t.Kind() == reflect.Slice {
			line = fmt.Sprintf("%s = r.GetBytes(length)", fieldName)
		} else {
			iterVar := fmt.Sprintf("_%d", iterLevel)
			initVar := fmt.Sprintf("%s, length_%d := 0, length", iterVar, iterLevel)
			line = fmt.Sprintf("for %s; %s<length_%d; %s++ { //%s of %s",
				initVar, iterVar, iterLevel, iterVar, t.Kind(), t.Elem().Kind())
			*lines = append(*lines, line)
			if t.Elem().Kind() == reflect.Interface || t.Elem().Kind() == reflect.Struct {
				line = fmt.Sprintf("%s[%s] = Rand%s(r)", fieldName, iterVar, typeName)
				*lines = append(*lines, line)
			} else {
				varName := fieldName + "[" + iterVar + "]"
				nl := ctx.genFieldRandLines(elemT, lines, varName, iterLevel+1)
				needLength = needLength || nl
			}
			line = "}"
		}
	case reflect.Interface:
		typePath := t.PkgPath() + "." + t.Name()
		alias, ok := ctx.ifcPath2Alias[typePath]
		if !ok {
			panic("Cannot find alias for:" + typePath)
		}
		line = fmt.Sprintf("%s = Rand%s(r) // interface_decode", fieldName, alias)
	case reflect.Ptr:
		panic("Should not reach here")
	case reflect.Struct:
		if _, ok := ctx.leafTypes[t.PkgPath()+"."+t.Name()]; ok {
			if isPtr {
				*lines = append(*lines, ctx.initPtrMember(fieldName, t))
				line = fmt.Sprintf("*(%s) = Rand%s(r)", fieldName, t.Name())
			} else {
				line = fmt.Sprintf("%s = Rand%s(r)", fieldName, t.Name())
			}
		} else {
			if isPtr {
				*lines = append(*lines, ctx.initPtrMember(fieldName, t))
			}
			nl := ctx.genStructRandLines(t, lines, fieldName, iterLevel)
			needLength = needLength || nl
			line = "// end of " + fieldName
		}
	default:
		panic(fmt.Sprintf("Unknown Kind %s", t.Kind()))
	}
	*lines = append(*lines, line)
	return needLength
}

func (ctx *prepareCtx) genStructRandLines(t reflect.Type, lines *[]string, varName string, iterLevel int) bool {
	needLength := false
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		nl := ctx.genFieldRandLines(field.Type, lines, varName+"."+field.Name, iterLevel)
		needLength = needLength || nl
	}
	return needLength
}

var headerLogics = `
func codonEncodeBool(w io.Writer, v bool) error {
	slice := []byte{0}
	if v {
		slice = []byte{1}
	}
	_, err := w.Write(slice)
	return err
}
func codonEncodeVarint(w io.Writer, v int64) error {
	var buf [10]byte
	n := binary.PutVarint(buf[:], v)
	_, err := w.Write(buf[0:n])
	return err
}
func codonEncodeInt8(w io.Writer, v int8) error {
	_, err := w.Write([]byte{byte(v)})
	return err
}
func codonEncodeInt16(w io.Writer, v int16) error {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], uint16(v))
	_, err := w.Write(buf[:])
	return err
}
func codonEncodeUvarint(w io.Writer, v uint64) error {
	var buf [10]byte
	n := binary.PutUvarint(buf[:], v)
	_, err := w.Write(buf[0:n])
	return err
}
func codonEncodeUint8(w io.Writer, v uint8) error {
	_, err := w.Write([]byte{byte(v)})
	return err
}
func codonEncodeUint16(w io.Writer, v uint16) error {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], v)
	_, err := w.Write(buf[:])
	return err
}
func codonEncodeFloat32(w io.Writer, v float32) error {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], math.Float32bits(v))
	_, err := w.Write(buf[:])
	return err
}
func codonEncodeFloat64(w io.Writer, v float64) error {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(v))
	_, err := w.Write(buf[:])
	return err
}
func codonEncodeByteSlice(w io.Writer, v []byte) error {
	err := codonEncodeVarint(w, int64(len(v)))
	if err != nil {
		return err
	}
	_, err = w.Write(v)
	return err
}
func codonEncodeString(w io.Writer, v string) error {
	return codonEncodeByteSlice(w, []byte(v))
}
func codonDecodeBool(bz []byte, n *int, err *error) bool {
	if len(bz) < 1 {
		*err = errors.New("Not enough bytes to read")
		return false
	}
	*n = 1
	*err = nil
	return bz[0]!=0
}
func codonDecodeInt(bz []byte, m *int, err *error) int {
	i, n := binary.Varint(bz)
	if n == 0 {
		// buf too small
		*err = errors.New("buffer too small")
	} else if n < 0 {
		// value larger than 64 bits (overflow)
		// and -n is the number of bytes read
		n = -n
		*err = errors.New("EOF decoding varint")
	}
	*m = n
	return int(i)
}
func codonDecodeInt8(bz []byte, n *int, err *error) int8 {
	if len(bz) < 1 {
		*err = errors.New("Not enough bytes to read")
		return 0
	}
	*err = nil
	*n = 1
	return int8(bz[0])
}
func codonDecodeInt16(bz []byte, n *int, err *error) int16 {
	if len(bz) < 2 {
		*err = errors.New("Not enough bytes to read")
		return 0
	}
	*n = 2
	*err = nil
	return int16(binary.LittleEndian.Uint16(bz[:2]))
}
func codonDecodeInt32(bz []byte, n *int, err *error) int32 {
	i := codonDecodeInt64(bz, n, err)
	return int32(i)
}
func codonDecodeInt64(bz []byte, m *int, err *error) int64 {
	i, n := binary.Varint(bz)
	if n == 0 {
		// buf too small
		*err = errors.New("buffer too small")
	} else if n < 0 {
		// value larger than 64 bits (overflow)
		// and -n is the number of bytes read
		n = -n
		*err = errors.New("EOF decoding varint")
	}
	*m = n
	*err = nil
	return int64(i)
}
func codonDecodeUint(bz []byte, n *int, err *error) uint {
	i := codonDecodeUint64(bz, n, err)
	return uint(i)
}
func codonDecodeUint8(bz []byte, n *int, err *error) uint8 {
	if len(bz) < 1 {
		*err = errors.New("Not enough bytes to read")
		return 0
	}
	*n = 1
	*err = nil
	return uint8(bz[0])
}
func codonDecodeUint16(bz []byte, n *int, err *error) uint16 {
	if len(bz) < 2 {
		*err = errors.New("Not enough bytes to read")
		return 0
	}
	*n = 2
	*err = nil
	return uint16(binary.LittleEndian.Uint16(bz[:2]))
}
func codonDecodeUint32(bz []byte, n *int, err *error) uint32 {
	i := codonDecodeUint64(bz, n, err)
	return uint32(i)
}
func codonDecodeUint64(bz []byte, m *int, err *error) uint64 {
	i, n := binary.Uvarint(bz)
	if n == 0 {
		// buf too small
		*err = errors.New("buffer too small")
	} else if n < 0 {
		// value larger than 64 bits (overflow)
		// and -n is the number of bytes read
		n = -n
		*err = errors.New("EOF decoding varint")
	}
	*m = n
	*err = nil
	return uint64(i)
}
func codonDecodeFloat64(bz []byte, n *int, err *error) float64 {
	if len(bz) < 8 {
		*err = errors.New("Not enough bytes to read")
		return 0
	}
	*n = 8
	*err = nil
	i := binary.LittleEndian.Uint64(bz[:8])
	return math.Float64frombits(i)
}
func codonDecodeFloat32(bz []byte, n *int, err *error) float32 {
	if len(bz) < 4 {
		*err = errors.New("Not enough bytes to read")
		return 0
	}
	*n = 4
	*err = nil
	i := binary.LittleEndian.Uint32(bz[:4])
	return math.Float32frombits(i)
}
func codonGetByteSlice(bz []byte, length int) ([]byte, int, error) {
	if len(bz) < length {
		return nil, 0, errors.New("Not enough bytes to read")
	}
	return bz[:length], length, nil
}
func codonDecodeString(bz []byte, n *int, err *error) string {
	var m int
	length := codonDecodeInt64(bz, &m, err)
	if *err != nil {
		return ""
	}
	var bs []byte
	var l int
	bs, l, *err = codonGetByteSlice(bz[m:], int(length))
	*n = m + l
	return string(bs)
}
`
