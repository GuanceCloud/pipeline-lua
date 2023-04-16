package point

import (
	"log"
	"reflect"

	"github.com/GuanceCloud/cliutils/point"
	lua "github.com/yuin/gopher-lua"
)

const luaPointTypeName = "point"

func registerPointType(L *lua.LState) {
	t := L.NewTypeMetatable(luaPointTypeName)
	L.SetGlobal("point", t)

	// static attributes
	L.SetField(t, "new", L.NewFunction(newPoint))
	L.SetField(t, "int", L.NewFunction(newInt))
	L.SetField(t, "bytes", L.NewFunction(newBytes))

	// methods
	L.SetField(t, "__index", L.SetFuncs(L.NewTable(), pointMethods))
}

var pointMethods = map[string]lua.LGFunction{
	"get": pointGetKeyValue,
	"set": pointSetKeyValues,
}

func newPoint(L *lua.LState) int {
	name := L.CheckString(1)
	kvs := L.CheckTable(2)
	opts := L.CheckTable(3)

	if opts != nil {
		opts.ForEach(nil) // TODO: extract options of the point
	}

	pt := point.NewPointV2([]byte(name), nil)

	kvs.ForEach(func(v1, v2 lua.LValue) {
		var (
			key   string
			value any
		)

		switch x := v1.(type) {
		case lua.LString:
			key = string(x)
		default:
			log.Printf("ignore key %q", v1)
			return // ignore the key
		}

		switch x := v2.(type) {
		case lua.LNumber:
			value = float64(x)
		case lua.LString:
			value = string(x)
		case lua.LBool:
			value = bool(x)
		default:
			log.Printf("ignore val %v", v2)
			return
		}

		log.Printf("add kv %q:%v", key, value)
		pt.Add([]byte(key), value)
	})

	ud := L.NewUserData()
	ud.Value = pt
	L.SetMetatable(ud, L.GetTypeMetatable(luaPointTypeName))
	L.Push(ud)
	return 1
}

func checkPoint(L *lua.LState) *point.Point {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*point.Point); ok {
		return v
	}
	L.ArgError(1, "point expected")
	return nil
}

func pointGetKeyValue(L *lua.LState) int {
	p := checkPoint(L)

	if L.GetTop() != 2 {
		L.ArgError(1, "expect key name")
		return 0
	}

	key := []byte(L.CheckString(2))

	log.Printf("try get key %q...", key)

	v := p.Get(key)

	switch x := v.(type) {
	case int64:
		L.Push(lua.LNumber(x))
	case bool:
		L.Push(lua.LBool(x))
	case float64:
		L.Push(lua.LNumber(x))
	case []byte:
		L.Push(lua.LString(x))
	case uint64:
		L.Push(lua.LNumber(x))
	case nil:
		L.Push(lua.LNil)
	default: // TODO: other types not support
		log.Printf("get value type %q, return nothing",
			reflect.TypeOf(v).String())
		return 0
	}
	return 1
}

func pointSetKeyValues(L *lua.LState) int {
	p := checkPoint(L)

	if L.GetTop() != 2 {
		L.ArgError(1, "expect key alue table")
		return 0
	}

	force := L.CheckBool(3)

	t := L.CheckTable(1)
	t.ForEach(func(k, v lua.LValue) {
		var (
			key []byte
			val any
		)

		switch x := k.(type) {
		case lua.LString:
			key = []byte(string(x))
		default:
			log.Printf("ignore key %v", k)
			return
		}

		switch x := v.(type) {
		case lua.LBool:
			val = bool(x)
		case lua.LNumber:
			val = float64(x)
		case lua.LString:
			val = string(x)
		default: // ignore other types

			log.Printf("ignore val %v", v)
			return
		}

		log.Printf("add kv: %q: %v", key, val)

		if force {
			p.MustAdd(key, val)
		} else {
			p.Add(key, val)
		}
	})

	return 0
}
