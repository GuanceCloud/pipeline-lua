package re

import (
	T "testing"

	"github.com/GuanceCloud/cliutils/point"
	luapt "github.com/GuanceCloud/pipeline-lua/pkg/point"
	"github.com/stretchr/testify/assert"
	libs "github.com/vadv/gopher-lua-libs"
	lua "github.com/yuin/gopher-lua"
)

func BenchmarkGrok(b *T.B) {

	cases := []struct {
		name, ls string
		pts      int
	}{
		{
			name: "10-pts-grok-year",
			pts:  10,
			ls: `
g = grok.new("%{YEAR:year}")

function main(pts)
	for i, pt in pairs(pts) do
		msg = pt:get("message")
		res = g:grok(msg)
		pt:set(res, false)
	end

	return pts
end`,
		},

		{
			name: "1000-pts-grok-year",
			pts:  1000,
			ls: `
g = grok.new("%{YEAR:year}")

function main(pts)
	for i, pt in pairs(pts) do
		msg = pt:get("message")
		res = g:grok(msg)
		pt:set(res, false)
	end

	return pts
end`,
		},
	}

	for _, c := range cases {
		L := lua.NewState()

		RegisterTypes(L)
		luapt.RegisterType(L)

		assert.NoError(b, L.DoString(c.ls))
		pts := point.NewRander(point.WithRandText(3)).Rand(c.pts)

		b.Run(c.name, func(b *T.B) {
			for i := 0; i < b.N; i++ {
				luapt.Call(L, pts)
			}
		})

		L.Close()
	}
}

func TestGrok(t *T.T) {
	t.Run("basic", func(t *T.T) {
		lscript := `
local log = require("log")
local info = log.new()

info:set_prefix("Lua > [INFO] ")
info:set_flags({date=true, time=true, longfile=true})

g = grok.new("%{YEAR:year}")

function main(pts)

	info:printf("%d pts...",  #pts)

	for i, pt in pairs(pts) do
		msg = pt:get("message")
		res = g:grok(msg)
		pt:set(res, false)
	end

	return pts
end`

		L := lua.NewState()
		defer L.Close()

		libs.Preload(L)
		RegisterTypes(L)
		luapt.RegisterType(L)

		assert.NoError(t, L.DoString(lscript))
		pts := point.NewRander(point.WithRandText(3)).Rand(2)

		res, err := luapt.Call(L, pts)
		assert.NoError(t, err)
		for _, pt := range res {
			t.Logf("%s", pt.LineProto())
		}
	})
}
