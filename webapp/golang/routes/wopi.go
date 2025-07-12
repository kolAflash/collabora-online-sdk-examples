package routes

import (
  "fmt"
  "io"
  "net/http"
  "regexp"
  "strings"
)


func Wopi(w http.ResponseWriter, r *http.Request) {
  rFile := regexp.MustCompile("^/wopi/files/([^/]+)/?(.*)$")
  match := rFile.FindStringSubmatch(r.RequestURI)
  fileId := match[1]
  contents := strings.HasPrefix(match[2], "contents")

  if fileId != "" {    
    if !contents {  // just fileId
      if r.Method == "GET" {
        getFile(w, r)
      }
    } else {  // contents
      if r.Method == "GET" {
        getFileContents(w, r)
      } else if r.Method == "POST" {
        postFileContents(w, r)
      }
    }
  }
}


func getFile(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, `{
  "BaseFileName": "test.txt",
  "Size": 11,
  "UserId": 1,
  "UserCanWrite": true,
  "EnableInsertRemoteImage": true
}`)
}


func getFileContents(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Hello world!")
}


func postFileContents(w http.ResponseWriter, r *http.Request) {
  //fmt.Println("wopi PutFile endpoint")  // debug
  bodyBytes, err := io.ReadAll(r.Body)
  if err == nil {
    //_ = bodyBytes
    fmt.Println(string(bodyBytes))  // debug
    w.WriteHeader(http.StatusOK)  // 200
  } else {
    fmt.Println("Not possible to get the file content.")
    w.WriteHeader(http.StatusNotFound)  // 404
  }
}
