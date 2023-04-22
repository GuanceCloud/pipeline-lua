package point

import (
	T "testing"

	"github.com/GuanceCloud/cliutils/point"
	"github.com/stretchr/testify/assert"
	libs "github.com/vadv/gopher-lua-libs"
	lua "github.com/yuin/gopher-lua"
)

func BenchmarkCallMain(b *T.B) {

	cases := []struct {
		name   string
		script string
	}{
		{
			name: "set-3-tags",
			script: `
function main(pts)
	tags = {t1='v1', t2='v2', t3='v3'}
	for i, pt in pairs(pts) do
		-- pt:set_tags({t1='v1', t2='v2'})
		pt:set_tags(tags)
	end
	return pts
end`,
		},

		{
			name: "set-3-tags-in-place",
			script: `
function main(pts)
	for i, pt in pairs(pts) do
		-- pt:set_tags({t1='v1', t2='v2'})
		pt:set_tags({t1='v1', t2='v2', t3='v3'})
	end
	return pts
end`,
		},

		{
			name: "pass-and-do-nothing",
			script: `
function main(pts)
	tags = {t1='v1', t2='v2', t3='v3'}
	for i, pt in pairs(pts) do
	end
	return pts
end`,
		},
	}

	for _, c := range cases {
		L := lua.NewState()

		libs.Preload(L)
		RegisterType(L)
		assert.NoError(b, L.DoString(c.script))
		pts := point.NewRander().Rand(10)

		b.Run(c.name, func(b *T.B) {
			for i := 0; i < b.N; i++ {
				Call(L, pts)
			}
		})

		L.Close()
	}
}

func TestCallMain(t *T.T) {
	t.Run("call-main", func(t *T.T) {
		lscript := `
local log = require("log")
local info = log.new()

info:set_prefix("Lua > [INFO] ")
info:set_flags({date=true, time=true})

function main(pts)
	for i, pt in pairs(pts) do
		info:printf("[%d] %s", i, pt:lineproto())
		pt:set_tags({t1='v1', t2='v2'})
	end
	return pts
end`
		L := lua.NewState()
		defer L.Close()

		libs.Preload(L)
		RegisterType(L)

		assert.NoError(t, L.DoString(lscript))
		pts := point.NewRander().Rand(2)

		res, err := Call(L, pts)
		assert.NoError(t, err)
		assert.Len(t, res, 2)

		assert.Equal(t, []byte("v1"), res[0].GetTag([]byte("t1")))
		assert.Equal(t, []byte("v2"), res[0].GetTag([]byte("t2")))

		assert.Equal(t, []byte("v1"), res[1].GetTag([]byte("t1")))
		assert.Equal(t, []byte("v2"), res[1].GetTag([]byte("t2")))

		for _, pt := range res {
			t.Logf(pt.LineProto())
		}

	})
}
