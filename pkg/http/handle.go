package http

import (
	"fmt"
	"io/ioutil"
	"klipper_power/pkg/plug"
	"log"
	"net/http"
	"time"
)

func HandleKlipper(uri string, conn *plug.Socket, during int) {
	for {
		log.Println("start get printer status")
		resp, err := http.Get(uri)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)

		}
		if parserContent(content) {
			log.Println("printer printing finished,try close printer power")
			if err = conn.Discover(); err != nil {
				panic(err)
			}
			if err := conn.PowerOff(); err != nil {
				panic(err)
			}
			log.Println("printer power closed")
		}
		time.Sleep(time.Duration(during) * time.Second)
	}

}

func parserContent(buffer []byte) bool {
	content := string(buffer)
	fmt.Println(content)
	return true
}
