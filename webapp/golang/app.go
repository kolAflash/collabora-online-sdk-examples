// clear; setsid --fork bash -c 'sleep 1; curl "http://127.0.0.1:3001/collaboraUrl?server=http://127.0.0.1:9980" >/dev/null 2>&1'; echo -e '\n\n\n'; kindunsafe go run .

package main

import (
  "fmt"
  "net/http"
  "strings"
)

import "CollaboraOnlineIntegrationExample/routes"


func main() {
  //const listenAddr = "localhost:3000";  // TODO won't work, IPv6 is needed!
  //const listenAddr = ":3000";  // important!  IPv6  // BAD, listens publicly
  const listenAddr = "[::1]:3000";  // important!  IPv6
  // TODO respect environment PORT variable
  fmt.Println("Listening on http://" + listenAddr + "/")

  http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
    fmt.Printf("%s: %s\n", r.Method, r.RequestURI)
    if strings.HasPrefix(r.RequestURI, "/collaboraUrl") {
      routes.CollaboraUrl(w, r)
    } else if strings.HasPrefix(r.RequestURI, "/wopi") {
      routes.Wopi(w, r)
    } else {
      routes.Index(w, r)
    }
  })
  http.ListenAndServe(listenAddr, nil)
}
