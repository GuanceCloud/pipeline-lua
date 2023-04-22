package point

import (
	"log"
	"time"

	"github.com/GuanceCloud/cliutils/point"
	lua "github.com/yuin/gopher-lua"
)

const luaPointTypeName = "point"

func RegisterType(L *lua.LState) {
	t := L.NewTypeMetatable(luaPointTypeName)
	L.SetGlobal("point", t)

	// static attributes
	L.SetField(t, "new", L.NewFunction(newPoint))
	L.SetField(t, "int", L.NewFunction(newInt))
	L.SetField(t, "uint", L.NewFunction(newUint))
	L.SetField(t, "bytes", L.NewFunction(newBytes))
	L.SetField(t, "tag", L.NewFunction(newTag))

	// methods
	L.SetField(t, "__index", L.SetFuncs(L.NewTable(), pointMethods))
}

var pointMethods = map[string]lua.LGFunction{
	"get":       pointGetKeyValue,
	"set":       pointSetKeyValues,
	"set_tags":  pointSetTags,
	"lineproto": pointLineProto,
}

func getTimeOpt(v lua.LValue) time.Time {
	switch x := v.(type) {
	case lua.LNumber:
		return time.Unix(0, int64(x))
	default:
		return time.Unix(0, 0)
	}
}

func newPoint(L *lua.LState) int {
	name := L.CheckString(1)
	kvs := L.CheckTable(2)
	opts := L.CheckTable(3)

	var ptopts []point.Option

	if opts != nil {
		opts.ForEach(func(v1, v2 lua.LValue) {
			switch k := v1.(type) {
			case lua.LString:
				switch k {
				case "time":
					ptopts = append(ptopts, point.WithTime(getTimeOpt(v2)))

				default:
					// TODO: add more options for point
					log.Printf("ignore option %s", k)
				}
			}
		})
	}

	var ptkvs point.KVs

	kvs.ForEach(func(v1, v2 lua.LValue) {
		var (
			key   string
			value any
		)

		isTag := false
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
		case *lua.LNilType:
			value = nil

		case *lua.LUserData:
			switch y := x.Value.(type) {
			case Int:
				value = int64(y)
			case Uint:
				value = uint64(y)
			case Bytes:
				value = []byte(y)
			case Tag:
				value = []byte(y)
				isTag = true
			}
		default:
			log.Printf("ignore val %v", v2)
			return
		}

		ptkvs = ptkvs.Add([]byte(key), value, isTag, false)
	})

	pt := point.NewPointV2([]byte(name), ptkvs, ptopts...)
	L.Push(luaPoint(L, pt))
	return 1
}

func luaPoint(L *lua.LState, pt *point.Point) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = pt
	L.SetMetatable(ud, L.GetTypeMetatable(luaPointTypeName))

	return ud
}

func checkPoint(L *lua.LState) *point.Point {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*point.Point); ok {
		return v
	}
	L.ArgError(1, "point expected")
	return nil
}
