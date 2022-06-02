package http

import (
	"klipper_power/pkg/plug"
	"net/http"
)

func StartServer(conn *plug.Socket) {
	http.Handle("/on", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if err := conn.Discover(); err != nil {
			panic(err)
		}
		if err := conn.PowerOn(); err != nil {
			panic(err)
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	}))

	if err := http.ListenAndServe(":2333", nil); err != nil {
		panic(err)
	}
}
