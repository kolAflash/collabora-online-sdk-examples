package routes

import (
  "encoding/xml"
  "fmt"
  "io"
  "net/http"
)


func Index(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" {
    http.FileServer(http.Dir("../html")).ServeHTTP(w, r)
  }
}


type xmlWopiDiscovery struct {
  XMLName xml.Name `xml:"wopi-discovery"`
  NetZoneAry []xmlNetZone `xml:"net-zone"`
}
type xmlNetZone struct { AppAry []xmlApp `xml:"app"` }
type xmlApp struct {
  Name string `xml:"name,attr"`
  ActionAry []xmlAction `xml:"action"`
}
type xmlAction struct { UrlSrc string `xml:"urlsrc,attr"` }

// GET: /collaboraUrl?server=http://127.0.0.1:9980
func CollaboraUrl(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" {
    var collaboraOnlineHost = r.URL.Query().Get("server")
    var clientResponse, _ = http.Get(collaboraOnlineHost + "/hosting/discovery")
    var xmlBytes, _ = io.ReadAll(clientResponse.Body)
    var wopiDiscovery xmlWopiDiscovery
    xml.Unmarshal(xmlBytes, &wopiDiscovery)
    // xpath: /wopi-discovery/net-zone/app[@name='text/plain']/action
    var appAry = wopiDiscovery.NetZoneAry[0].AppAry
    var onlineUrl string
    for i := len(appAry)-1; i > -1; i-- {
      var appElem = appAry[i]
      if appElem.Name == "text/plain" {
        onlineUrl = appElem.ActionAry[0].UrlSrc
        i = -1
      }
    }
    fmt.Fprint(w, `{"url":"` + onlineUrl + `","token":"test"}`)
  }
}




// fmt.Printf("             r: %s\n", r)
// fmt.Printf("       Headers: %s\n", r.Header)
// fmt.Printf("          Host: %s\n", r.Host)
// fmt.Printf("           URL: %s\n", r.URL)
// fmt.Printf("        Method: %s\n", r.Method)
// fmt.Printf("    RequestURI: %s\n", r.RequestURI)
// fmt.Printf("      URL.Path: %s\n", r.URL.Path)
// fmt.Printf("   URL.RawPath: %s\n", r.URL.RawPath)
// fmt.Printf("  URL.RawQuery: %s\n", r.URL.RawQuery)
// fmt.Printf("URL.RequestURI: %s\n", r.URL.RequestURI())
// fmt.Printf("    URL.Scheme: %s\n", r.URL.Scheme)
// fmt.Printf("      URL.Host: %s\n", r.URL.Host)
// fmt.Printf("  URL.Hostname: %s\n", r.URL.Hostname())
// fmt.Printf("     URL.Query: %s\n", r.URL.Query())

// var dec = xml.NewDecoder(strings.NewReader(string(xmlBytes)))
// var stack []string
// var objStack []xml.StartElement
// var result string
// for {
//   var tok, xmlErr = dec.Token()
//   if xmlErr != nil { break }
//   switch se := tok.(type) {
//   case xml.StartElement:
//     stack = append(stack, se.Name.Local)
//     objStack = append(objStack, se)
//   case xml.EndElement:
//     stack = stack[:(len(stack) - 1)]  // pop
//     objStack = objStack[:(len(objStack) - 1)]  // pop
//   }
//   if fmt.Sprintf("%s", stack) == "[wopi-discovery net-zone app action]" {
//     fmt.Printf("%s\n\n", objStack)
//     for _, appAttr := range objStack[2].Attr {
//       if appAttr.Name.Local == "name" && appAttr.Value == "text/plain" {
//         for _, actionAttr := range objStack[3].Attr {
//           if actionAttr.Name.Local == "urlsrc" {
//             fmt.Println(actionAttr.Value)
//             result = actionAttr.Value
//             break
//           }
//         }
//         break
//       }
//     }
//   }
//   if len(result) != 0 { break }
// }
