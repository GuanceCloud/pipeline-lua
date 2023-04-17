package point

import (
	lua "github.com/yuin/gopher-lua"
)

const (
	luaBytesTypeName = "bytes"
	luaIntTypeName   = "int"
	luaUintTypeName  = "uint"
	luaTagTypeName   = "tag"
)

type Bytes []byte
type Int int64
type Uint uint64
type Tag string

func newBytes(L *lua.LState) int {
	s := L.CheckString(1)

	ud := L.NewUserData()
	ud.Value = Bytes(string(s))
	L.SetMetatable(ud, L.GetTypeMetatable(luaBytesTypeName))
	L.Push(ud)
	return 1
}

func newInt(L *lua.LState) int {
	i := L.CheckNumber(1)

	ud := L.NewUserData()
	ud.Value = Int(i)
	L.SetMetatable(ud, L.GetTypeMetatable(luaIntTypeName))
	L.Push(ud)
	return 1
}

func newUint(L *lua.LState) int {
	i := L.CheckNumber(1)

	ud := L.NewUserData()
	ud.Value = Uint(i)
	L.SetMetatable(ud, L.GetTypeMetatable(luaUintTypeName))
	L.Push(ud)
	return 1
}

func newTag(L *lua.LState) int {
	s := L.CheckString(1)

	ud := L.NewUserData()
	ud.Value = Tag(s)
	L.SetMetatable(ud, L.GetTypeMetatable(luaTagTypeName))
	L.Push(ud)
	return 1
}
