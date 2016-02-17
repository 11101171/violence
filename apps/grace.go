package main

import (
	"net/http"
	"os"
	"strconv"
)

//****************************************************/
//Copyright(c) 2015 Tencent, all rights reserved
// File        : apps/grace.go
// Author      : ningzhong.zeng
// Revision    : 2015-12-25 16:13:01
// Description :
//****************************************************/

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("WORLD!"))
	w.Write([]byte("ospid:" + strconv.Itoa(os.Getpid())))
}

func main() {
	/*  mux := http.NewServeMux() */
	// mux.HandleFunc("/hello", handler)

	// err := grace.ListenAndServe("localhost:9090", mux)
	// if err != nil {
	// log.Println(err)
	// }
	// log.Println("Server on 9090 stopped")
	/* os.Exit(0) */
}
