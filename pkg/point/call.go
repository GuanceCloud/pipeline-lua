package point

import (
	"fmt"
	"reflect"

	"github.com/GuanceCloud/cliutils/point"
	lua "github.com/yuin/gopher-lua"
)

const (
	Entry = "main"
)

// Call is the entry to pass points to lua.
func Call(L *lua.LState, pts []*point.Point) ([]*point.Point, error) {
	t := L.NewTable()
	for _, pt := range pts {
		t.Append(luaPoint(L, pt))
	}

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal(Entry),
		NRet:    1,
		Protect: true,
	}, t); err != nil {
		return nil, err
	}

	ret := L.Get(-1)
	L.Pop(1)

	switch x := ret.(type) {
	case *lua.LTable:
		var res []*point.Point
		x.ForEach(func(v1, v2 lua.LValue) {
			switch ud := v2.(type) {
			case *lua.LUserData:
				res = append(res, ud.Value.(*point.Point))
			}
		})

		return res, nil

	default:
		return nil, fmt.Errorf("unexpected result type: %s", reflect.TypeOf(ret).String())
	}
}
