package httpserver

import (
	//"encoding/json"
	"bytes"
	"fmt"
	"log"
	"net/http"
)

func Start() {
	go func() {
		http.HandleFunc("/push", test)
		log.Fatal(http.ListenAndServe(":8082", nil))
	}()

}

func test(rw http.ResponseWriter, req *http.Request) {
	//decoder := json.NewDecoder(req.Body)
	// var t test_struct
	// err = decoder.Decode(&t)
	// if err != nil {
	// 	panic()
	// }

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	body := buf.String()
	fmt.Println(body)
	//log.Println(req.Body)
}
