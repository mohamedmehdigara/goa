package codegen

import (
	"testing"

	"goa.design/goa/design"
)

var (
	SimpleObj     = require(object("a", design.String, "b", design.Int), "a")
	RequiredObj   = require(object("a", design.String, "b", design.Int), "a", "b")
	DefaultObj    = defaulta(require(object("a", SimpleArray.Type, "b", design.Int), "a"), "a", []string{"foo", "bar"}, "b", 42)
	SuperObj      = require(object("a", design.String, "b", design.Int, "c", design.Boolean), "a")
	SimpleArray   = array(design.String)
	SimpleMap     = mapa(design.String, design.Int)
	ArrayObj      = object("a", design.String, "b", SimpleArray.Type)
	CompositeObj  = defaulta(require(object("aa", SimpleArray.Type, "bb", SimpleObj.Type), "bb"), "aa", []string{"foo", "bar"})
	ObjArray      = array(RequiredObj.Type)
	ObjMap        = mapa(design.String, SimpleObj.Type)
	UserType      = object("ut", &design.UserTypeExpr{TypeName: "User", AttributeExpr: SimpleObj})
	ArrayUserType = array(&design.UserTypeExpr{TypeName: "User", AttributeExpr: RequiredObj})
)

func TestGoTypeTransform(t *testing.T) {
	var (
		sourceVar = "source"
		targetVar = "target"
	)
	cases := []struct {
		Name           string
		Source, Target *design.AttributeExpr
		Marshal        bool
		TargetPkg      string

		Code string
	}{
		{"simple-unmarshal", SimpleObj, SimpleObj, true, "", objUnmarshalCode},
		{"required-unmarshal", SimpleObj, RequiredObj, true, "", requiredUnmarshalCode},
		{"default-unmarshal", DefaultObj, DefaultObj, true, "", defaultUnmarshalCode},

		{"simple-marshal", SimpleObj, SimpleObj, false, "", objCode},
		{"required-marshal", RequiredObj, RequiredObj, false, "", requiredCode},
		{"default-marshal", DefaultObj, DefaultObj, false, "", defaultCode},

		// non match field ignore
		{"super-unmarshal", SuperObj, SimpleObj, true, "", objUnmarshalCode},
		{"super-marshal", SuperObj, SimpleObj, false, "", objCode},
		{"super-unmarshal-r", SimpleObj, SuperObj, true, "", objUnmarshalCode},
		{"super-marshal-r", SimpleObj, SuperObj, false, "", objCode},

		// simple array and map
		{"array-unmarshal", SimpleArray, SimpleArray, true, "", arrayCode},
		{"map-unmarshal", SimpleMap, SimpleMap, true, "", mapCode},
		{"object-array-unmarshal", ArrayObj, ArrayObj, true, "", arrayObjUnmarshalCode},

		{"array-marshal", SimpleArray, SimpleArray, false, "", arrayCode},
		{"map-marshal", SimpleMap, SimpleMap, false, "", mapCode},
		{"object-array-marshal", ArrayObj, ArrayObj, false, "", arrayObjCode},

		// composite data structures
		{"composite-unmarshal", CompositeObj, CompositeObj, true, "", compUnmarshalCode},
		{"composite-marshal", CompositeObj, CompositeObj, false, "", compCode},

		// object in arrays and maps
		{"object-array-unmarshal", ObjArray, ObjArray, true, "", objArrayCode},
		{"object-map-unmarshal", ObjMap, ObjMap, true, "", objMapCode},
		{"user-type-unmarshal", UserType, UserType, true, "", userTypeUnmarshalCode},
		{"array-user-type-unmarshal", ArrayUserType, ArrayUserType, true, "", arrayUserTypeUnmarshalCode},

		{"object-array-marshal", ObjArray, ObjArray, false, "", objArrayCode},
		{"object-map-marshal", ObjMap, ObjMap, false, "", objMapCode},
		{"user-type-marshal", UserType, UserType, false, "", userTypeCode},
		{"array-user-type-marshal", ArrayUserType, ArrayUserType, false, "", arrayUserTypeCode},

		// package handling
		{"target-package-unmarshal", ArrayUserType, ArrayUserType, true, "tpkg", objTargetPkgUnmarshalCode},
		{"target-package-marshal", ArrayUserType, ArrayUserType, false, "tpkg", objTargetPkgCode},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			src := &design.UserTypeExpr{TypeName: "SourceType", AttributeExpr: c.Source}
			tgt := &design.UserTypeExpr{TypeName: "TargetType", AttributeExpr: c.Target}
			code, _, err := GoTypeTransform(src, tgt, sourceVar, targetVar, "", c.TargetPkg, c.Marshal, NewNameScope())
			if err != nil {
				t.Fatal(err)
			}
			code = FormatTestCode(t, "package foo\nfunc transform(){\n"+code+"}")
			if code != c.Code {
				t.Errorf("invalid code, got:\n%s\ngot vs. expected:\n%s", code, Diff(t, code, c.Code))
			}
		})
	}
}

func require(att *design.AttributeExpr, names ...string) *design.AttributeExpr {
	att.Validation = &design.ValidationExpr{Required: names}
	return att
}

func defaulta(att *design.AttributeExpr, vals ...interface{}) *design.AttributeExpr {
	obj := att.Type.(*design.Object)
	for i := 0; i < len(vals); i += 2 {
		name := vals[i].(string)
		obj.Attribute(name).DefaultValue = vals[i+1]
	}
	return att
}

func object(params ...interface{}) *design.AttributeExpr {
	obj := design.Object{}
	for i := 0; i < len(params); i += 2 {
		name := params[i].(string)
		typ := params[i+1].(design.DataType)
		obj = append(obj, &design.NamedAttributeExpr{Name: name, Attribute: &design.AttributeExpr{Type: typ}})
	}
	return &design.AttributeExpr{Type: &obj}
}

func array(dt design.DataType) *design.AttributeExpr {
	elem := &design.AttributeExpr{Type: dt}
	return &design.AttributeExpr{Type: &design.Array{ElemType: elem}}
}

func mapa(keyt, elemt design.DataType) *design.AttributeExpr {
	key := &design.AttributeExpr{Type: keyt}
	elem := &design.AttributeExpr{Type: elemt}
	return &design.AttributeExpr{Type: &design.Map{KeyType: key, ElemType: elem}}
}

const objUnmarshalCode = `func transform() {
	target := &TargetType{
		A: *source.A,
		B: source.B,
	}
}
`

const requiredUnmarshalCode = `func transform() {
	target := &TargetType{
		A: *source.A,
	}
	if source.B != nil {
		target.B = *source.B
	}
}
`

const defaultUnmarshalCode = `func transform() {
	target := &TargetType{}
	if source.B != nil {
		target.B = *source.B
	}
	target.A = make([]string, len(source.A))
	for j, val := range source.A {
		target.A[j] = val
	}
	if source.A == nil {
		target.A = []string{"foo", "bar"}
	}
	if source.B == nil {
		target.B = 42
	}
}
`

const objCode = `func transform() {
	target := &TargetType{
		A: source.A,
		B: source.B,
	}
}
`

const requiredCode = `func transform() {
	target := &TargetType{
		A: source.A,
		B: source.B,
	}
}
`

const defaultCode = `func transform() {
	target := &TargetType{
		B: source.B,
	}
	if source.A != nil {
		target.A = make([]string, len(source.A))
		for j, val := range source.A {
			target.A[j] = val
		}
	}
	if source.A == nil {
		target.A = []string{"foo", "bar"}
	}
}
`

const objDefaultPointersCode = `func transform() {
	target := &TargetType{
		A: *source.A,
		B: source.B,
	}
	if source.B == nil {
		tmp := 42
		target.B = &tmp
	}
}
`

const arrayUnmarshalCode = `func transform() {
	target := make([]string, len(source))
	for i, val := range source {
		target[i] = val
	}
}
`

const arrayCode = `func transform() {
	target := make([]string, len(source))
	for i, val := range source {
		target[i] = val
	}
}
`
const arrayObjUnmarshalCode = `func transform() {
	target := &TargetType{
		A: source.A,
	}
	if source.B != nil {
		target.B = make([]string, len(source.B))
		for j, val := range source.B {
			target.B[j] = val
		}
	}
}
`

const arrayObjCode = `func transform() {
	target := &TargetType{
		A: source.A,
	}
	if source.B != nil {
		target.B = make([]string, len(source.B))
		for j, val := range source.B {
			target.B[j] = val
		}
	}
}
`

const mapCode = `func transform() {
	target := make(map[string]int, len(source))
	for key, val := range source {
		tk := key
		tv := val
		target[tk] = tv
	}
}
`

const compUnmarshalCode = `func transform() {
	target := &TargetType{}
	if source.Aa != nil {
		target.Aa = make([]string, len(source.Aa))
		for j, val := range source.Aa {
			target.Aa[j] = val
		}
	}
	if source.Aa == nil {
		target.Aa = []string{"foo", "bar"}
	}
	target.Bb = &struct {
		A *string
		B *int
	}{
		A: source.Bb.A,
		B: source.Bb.B,
	}
}
`

const compCode = `func transform() {
	target := &TargetType{}
	if source.Aa != nil {
		target.Aa = make([]string, len(source.Aa))
		for j, val := range source.Aa {
			target.Aa[j] = val
		}
	}
	if source.Aa == nil {
		target.Aa = []string{"foo", "bar"}
	}
	if source.Bb != nil {
		target.Bb = &struct {
			A *string
			B *int
		}{
			A: source.Bb.A,
			B: source.Bb.B,
		}
	}
}
`

const compDefaultsPointersCode = `func transform() {
	target := &TargetType{}
	if source.Aa != nil {
		target.Aa = *source.Aa
	}
	if source.Aa == nil {
		target.Aa = "default"
	}
	target.Bb = &struct {
		A *string
		B *int
	}{
		A: source.Bb.A,
		B: source.Bb.B,
	}
}
`

const objArrayCode = `func transform() {
	target := make([]struct {
		A *string
		B *int
	}, len(source))
	for i, val := range source {
		target[i] = &struct {
			A *string
			B *int
		}{
			A: val.A,
			B: val.B,
		}
	}
}
`

const objMapCode = `func transform() {
	target := make(map[string]struct {
		A *string
		B *int
	}, len(source))
	for key, val := range source {
		tk := key
		tv := &struct {
			A *string
			B *int
		}{
			A: val.A,
			B: val.B,
		}
		target[tk] = tv
	}
}
`

const userTypeUnmarshalCode = `func transform() {
	target := &TargetType{}
	if source.Ut != nil {
		target.Ut = unmarshalUserToUser(source.Ut)
	}
}
`

const userTypeCode = `func transform() {
	target := &TargetType{}
	if source.Ut != nil {
		target.Ut = marshalUserToUser(source.Ut)
	}
}
`

const arrayUserTypeUnmarshalCode = `func transform() {
	target := make([]*User, len(source))
	for i, val := range source {
		target[i] = &User{
			A: *val.A,
			B: *val.B,
		}
	}
}
`

const arrayUserTypeCode = `func transform() {
	target := make([]*User, len(source))
	for i, val := range source {
		target[i] = &User{
			A: val.A,
			B: val.B,
		}
	}
}
`

const objTargetPkgUnmarshalCode = `func transform() {
	target := make([]*tpkg.User, len(source))
	for i, val := range source {
		target[i] = &tpkg.User{
			A: *val.A,
			B: *val.B,
		}
	}
}
`

const objTargetPkgCode = `func transform() {
	target := make([]*tpkg.User, len(source))
	for i, val := range source {
		target[i] = &tpkg.User{
			A: val.A,
			B: val.B,
		}
	}
}
`