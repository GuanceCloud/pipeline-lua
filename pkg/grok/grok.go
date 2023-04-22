package re

import (
	"fmt"

	"github.com/GuanceCloud/grok"
	lua "github.com/yuin/gopher-lua"
)

const (
	luaGrokTypeName = "grok"
)

func RegisterTypes(L *lua.LState) {
	t := L.NewTypeMetatable(luaGrokTypeName)
	L.SetGlobal("grok", t)

	// static attributes
	L.SetField(t, "new", L.NewFunction(newGrok))

	// methods
	L.SetField(t, "__index", L.SetFuncs(L.NewTable(), methods))
}

var methods = map[string]lua.LGFunction{
	"grok": doGrok,
}

type patterns struct {
	x map[string]*grok.GrokPattern
}

func (p patterns) GetPattern(k string) (*grok.GrokPattern, bool) {
	if p.x == nil {
		return nil, false
	}

	ptn, ok := p.x[k]

	return ptn, ok
}

func (p patterns) SetPattern(k string, ptn *grok.GrokPattern) {
	if p.x == nil {
		p.x = map[string]*grok.GrokPattern{}
	}

	p.x[k] = ptn
}

type luaGrok struct {
	ptn string
	g   *grok.GrokRegexp
}

func newGrok(L *lua.LState) int {
	pattern := L.CheckString(1)

	de, errs := grok.DenormalizePatternsFromMap(grok.CopyDefalutPatterns())
	if len(errs) != 0 {
		L.ArgError(1, fmt.Sprintf("grok load pattern failed: %+#v", errs))
	}

	g, err := grok.CompilePattern(pattern, patterns{x: de})
	if err != nil {
		L.ArgError(1, fmt.Sprintf("grok %q compiled error: %s", pattern, err))
		return 0
	}

	L.Push(udGrok(L, &luaGrok{ptn: pattern, g: g}))
	return 1
}

func udGrok(L *lua.LState, g *luaGrok) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = g
	L.SetMetatable(ud, L.GetTypeMetatable(luaGrokTypeName))
	return ud
}

func doGrok(L *lua.LState) int {
	g := checkGrok(L)

	if L.GetTop() != 2 {
		L.ArgError(1, "expect text/message to grok")
		return 0
	}

	msg := L.CheckString(2)

	res, _, err := g.g.RunWithTypeInfo(msg, true)
	if err != nil {
		L.ArgError(1, fmt.Sprintf("grok run error: %s", err))
		return 0
	}

	t := L.NewTable()
	for k, v := range res {
		switch x := v.(type) {
		case int64:
			t.RawSet(lua.LString(k), lua.LNumber(x))
		case float64:
			t.RawSet(lua.LString(k), lua.LNumber(x))
		case bool:
			t.RawSet(lua.LString(k), lua.LBool(x))
		case string:
			t.RawSet(lua.LString(k), lua.LString(x))
		}
	}

	L.Push(t)
	return 1
}

func checkGrok(L *lua.LState) *luaGrok {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*luaGrok); ok {
		return v
	}
	L.ArgError(1, "point expected")
	return nil
}
