package point

import lua "github.com/yuin/gopher-lua"

const (
	luaBytesTypeName = "bytes"
	luaIntTypeName   = "int"
)

func newBytes(L *lua.LState) int {
	s := L.CheckString(1)

	ud := L.NewUserData()
	ud.Value = []byte(string(s))
	L.SetMetatable(ud, L.GetTypeMetatable(luaBytesTypeName))
	L.Push(ud)
	return 1
}

func newInt(L *lua.LState) int {
	i := L.CheckNumber(1)

	ud := L.NewUserData()
	ud.Value = int64(i)
	L.SetMetatable(ud, L.GetTypeMetatable(luaIntTypeName))
	L.Push(ud)
	return 1
}
