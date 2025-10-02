package template

import _ "embed"

// LuaTemplate contains the embedded Lua template for code generation
//
//go:embed lua.gotmpl
var LuaTemplate string
