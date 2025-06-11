package app

import (
	"fmt"
	"gontainer/internal/controller/value/handler"
	"gontainer/internal/kgsp/reader"
	"gontainer/internal/kgsp/spvalue"
	"gontainer/internal/kgsp/writer"
	"net"
	"strings"
)

func Run() error {
	fmt.Println("Listening on port :6379")

	// Create a new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		return err
	}

	var handlers = map[string]func([]spvalue.Value) spvalue.Value{
		"P": handler.HandlePing,
	}

	// Listen for connections
	conn, err := l.Accept()
	if err != nil {
		return err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()

	for {
		resp := reader.NewSPReader(conn)
		value, err := resp.Read()
		if err != nil {
			continue
		}

		fmt.Println(string(spvalue.Value{Typ: "string", Str: "hi"}.Marshal()))

		if value.Typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.Array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		w := writer.NewWriter(conn)

		h, ok := handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			err := w.Write(spvalue.Value{Typ: "string", Str: ""})
			if err != nil {
				return err
			}
			continue
		}

		result := h(args)
		err = w.Write(result)
		if err != nil {
			return err
		}
	}
}
