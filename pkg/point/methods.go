package point

import (
	"log"
	"reflect"

	"github.com/GuanceCloud/cliutils/point"
	lua "github.com/yuin/gopher-lua"
)

func pointLineProto(L *lua.LState) int {
	p := checkPoint(L)

	lp := p.LineProto()
	L.Push(lua.LString(lp))
	return 1
}

func pointSetTags(L *lua.LState) int {
	p := checkPoint(L)

	if L.GetTop() != 2 {
		L.ArgError(1, "expect tag list")
		return 0
	}

	t := L.CheckTable(2)
	tags := map[string]string{}
	t.ForEach(func(k, v lua.LValue) {
		switch x := k.(type) {
		case lua.LString:
			switch y := v.(type) {
			case lua.LString:
				tags[string(x)] = string(y)
			}

			// others ignored
		}
	})

	for k, v := range tags {
		p.AddTag([]byte(k), []byte(v))
	}

	return 0
}

func pointGetKeyValue(L *lua.LState) int {
	p := checkPoint(L)

	if L.GetTop() != 2 {
		L.ArgError(1, "expect key name")
		return 0
	}

	key := []byte(L.CheckString(2))

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

	if L.GetTop() != 3 {
		L.ArgError(1, "expect key value table")
		return 0
	}

	t := L.CheckTable(2)
	force := L.CheckBool(3)

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

		isTag := false

		switch x := v.(type) {
		case lua.LBool:
			val = bool(x)
		case lua.LNumber:
			val = float64(x)
		case lua.LString:
			val = string(x)
		case *lua.LUserData:

			switch y := x.Value.(type) {
			case Int:
				val = int64(y)
			case Uint:
				val = uint64(y)
			case Bytes:
				val = []byte(y)
			case Tag:
				val = []byte(y)
				isTag = true
			}

		default: // ignore other types

			log.Printf("ignore val %v", v)
			return
		}

		kv := point.NewKV(key, val, point.WithKVTagSet(isTag))

		if force {
			p.MustAddKV(kv)
		} else {
			p.AddKV(kv)
		}
	})

	return 0
}
