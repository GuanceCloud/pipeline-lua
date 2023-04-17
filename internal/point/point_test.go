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

kvs  = { f1 = 1.1, f2 = "abc" }
p = point.new("test", kvs, {time=0})  -- with timestamp 0

info:printf("f1: %f", p:get("f1"))
info:printf("f2: %s", p:get("f2"))
info:printf("lp: %s", p:lineproto())
assert(p:lineproto() == 'test f1=1.1,f2="abc" 0')
`))
	})

	t.Run("set-point-tag", func(t *T.T) {
		L := lua.NewState()
		defer L.Close()

		libs.Preload(L)

		registerPointType(L)
		assert.NoError(t, L.DoString(`
local log = require("log")
local info = log.new()

info:set_prefix("Lua > [INFO] ")
info:set_flags({data=true, time=true})

kvs  = { f1 = 1, f2 = "abc", }
p = point.new("test", kvs, {time=1}) 

p:set_tags({t1 = 'v1', t2 = "v2"})

info:printf("f1: %f", p:get("f1"))
info:printf("f2: %s", p:get("f2"))
info:printf("lp: %s", p:lineproto())
assert(p:lineproto() == 'test,t1=v1,t2=v2 f1=1,f2="abc" 1')
`))
	})

	t.Run("new-point-with-tag", func(t *T.T) {
		L := lua.NewState()
		defer L.Close()

		libs.Preload(L)

		registerPointType(L)
		assert.NoError(t, L.DoString(`
local log = require("log")
local info = log.new()

info:set_prefix("Lua > [INFO] ")
info:set_flags({data=true, time=true})

kvs  = { f1 = 1, f2 = "abc", t1= point.tag("v1") }
p = point.new("test", kvs, {time=1}) 

info:printf("f1: %f", p:get("f1"))
info:printf("f2: %s", p:get("f2"))
info:printf("lp: %s", p:lineproto())
assert(p:lineproto() == 'test,t1=v1 f1=1,f2="abc" 1')
`))
	})

	t.Run("new-point-with-time", func(t *T.T) {
		L := lua.NewState()
		defer L.Close()

		libs.Preload(L)

		registerPointType(L)
		assert.NoError(t, L.DoString(`
local log = require("log")
local info = log.new()

info:set_prefix("Lua > [INFO] ")
info:set_flags({data=true, time=true})

kvs  = { f1 = 1, f2 = "abc", }
p = point.new("test", kvs, {time = 123}) 

info:printf("f1: %f", p:get("f1"))
info:printf("f2: %s", p:get("f2"))
info:printf("lp: %s", p:lineproto())
assert(p:lineproto() == 'test f1=1,f2="abc" 123')
`))
	})

	t.Run("set-point-value-as-int", func(t *T.T) {
		L := lua.NewState()
		defer L.Close()

		libs.Preload(L)
		registerPointType(L)

		assert.NoError(t, L.DoString(`
local log = require("log")
local info = log.new()
info:set_prefix("lua > [INFO]")
info:set_flags({data=true, time=true})

kvs = { f1=point.int(1), f2 = point.bytes('abc'), f3="xyz" }
p = point.new("test", kvs, {time = 456})

assert(p:lineproto() == 'test f1=1i,f2="abc",f3="xyz" 456')
		`))
	})

	t.Run("set-point-kvs", func(t *T.T) {
		L := lua.NewState()
		defer L.Close()

		libs.Preload(L)
		registerPointType(L)

		assert.NoError(t, L.DoString(`
local log = require("log")
local info = log.new()
info:set_prefix("lua > [INFO]")
info:set_flags({data=true, time=true})

p = point.new("test", {}, {time = 456})
p:set({
  f1 = 123, -- float
  f2 = point.int(321), -- int64
  f3 = point.uint(456), -- uint64
  t1 = point.tag("v1"), -- tag
}, false) 

info:printf("lp: %s", p:lineproto())

-- 456u uint also as int(456i) in line-protocol(inflxudb:v1)
assert(p:lineproto() == 'test,t1=v1 f1=123,f2=321i,f3=456i 456') 
		`))
	})
}
