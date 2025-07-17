// clear; setsid --fork bash -c 'sleep 1; curl "http://127.0.0.1:3001/collaboraUrl?server=http://127.0.0.1:9980" >/dev/null 2>&1'; echo -e '\n\n\n'; kindunsafe go run .

package main

import (
  "fmt"
  "net/http"
  "strings"
)

import "CollaboraOnlineIntegrationExample/routes"


// TODO: print all errors from unexpected behaviour!
func main() {
  //const listenAddr = "localhost:3000";  // TODO won't work, COOL expects IPv6 host!
  const listenAddr = ":3000";  // important!  IPv6+IPv4 needed  // BAD, listens publicly
  //const listenAddr = "[::1]:3000";  // TODO: won't work, Firefox doesn't allow embedding v4 in v6 localhost
  // TODO respect environment PORT variable
  fmt.Println("Listening on http://" + listenAddr + "/")

  http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
    fmt.Printf("%s: %s\n", r.Method, r.RequestURI)
    if strings.HasPrefix(r.RequestURI, "/wopi") {
      routes.Wopi(w, r)
    } else if strings.HasPrefix(r.RequestURI, "/collaboraUrl") {
      routes.CollaboraUrl(w, r)
    } else {
      routes.Index(w, r)
    }
  })
  http.ListenAndServe(listenAddr, nil)
}
