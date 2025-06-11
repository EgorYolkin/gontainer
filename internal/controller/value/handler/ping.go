package handler

import (
	"gontainer/internal/kgsp/spvalue"
)

func HandlePing(args []spvalue.Value) spvalue.Value {
	return spvalue.Value{Typ: "string", Str: "PONG"}
}
