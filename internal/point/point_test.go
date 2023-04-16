package point

import (
	T "testing"

	"github.com/stretchr/testify/assert"
	libs "github.com/vadv/gopher-lua-libs"
	lua "github.com/yuin/gopher-lua"
)

func TestPoint(t *T.T) {
	t.Run("new-point", func(t *T.T) {
		L := lua.NewState()
		defer L.Close()

		libs.Preload(L)

		registerPointType(L)
		assert.NoError(t, L.DoString(`
local log = require("log")
local info = log.new()

info:set_prefix("Lua > [INFO] ")
info:set_flags({data=true, time=true})

kvs  = { f1 = 1, f2 = "abc" }
p = point.new("test", kvs, {}) 

info:printf("f1: %f", p:get("f1"))
info:printf("f2: %s", p:get("f2"))
`))
	})
}
