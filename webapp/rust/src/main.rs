use std::{
  fs,
  io::{BufReader, prelude::*},
  net::TcpListener
};

use curl::easy::Easy;

use sxd_document::parser;
use sxd_xpath::evaluate_xpath;




// TODO: print all errors from unexpected behaviour!
fn main() {
  let port = std::env::var("PORT").ok()
    .and_then(|s| s.parse::<u16>().ok()).unwrap_or(3000);

  // TODO won't work, COOL expects IPv6 host!
  //let addr = std::net::SocketAddr::from((std::net::Ipv4Addr::LOCALHOST, port));
  // important!  IPv6+IPv4 needed  // BAD, listens publicly
  let addr_string = format!("[::0]:{}", port);
  let addr = addr_string.as_str();
  // TODO: won't work, Firefox doesn't allow embedding v4 in v6 localhost
  //let addr = std::net::SocketAddr::from((std::net::Ipv6Addr::LOCALHOST, port));

  let listener = TcpListener::bind(addr).unwrap();
  println!("Listening on http://{}/", addr);
  // TODO: no HTTPS support yet


  for w_stream in listener.incoming() {
    println!("connection received: {:?}", w_stream);
    if let Ok(mut stream) = w_stream {
      let mut buf_reader = BufReader::new(&stream);
      let mut lines = buf_reader.by_ref().lines();
      if let Some(Ok(method_line)) = lines.next() {
        println!("HTTP method line: {:?}", method_line);  // example: GET /index.html HTTP/1.1
        let mut method_args = method_line.split(" ");
        let method = method_args.next().unwrap_or("");
        let url = method_args.next().unwrap_or("");
        let mut response: Option<String> = None;
      
        if response.is_none() && url.starts_with("/wopi/files/") {
          let sub_url = &url[12..url.len()];
          let file_id_end = sub_url.find('/').unwrap_or(sub_url.len());
          let file_id = &sub_url[0..file_id_end];
          let contents = &sub_url[file_id_end..sub_url.len()].starts_with("/contents?");
          if file_id != "" {
            if !contents {  // just file_id
              if method == "GET" { response = Some(GET_FILE_JSON.to_string()) }
            } else {  // contents
              if method == "GET" {
                response = Some("Hello world!".to_string());
              } else if method == "POST" {
                response = Some("".to_string());
              }
            }
          }
        }
      
        if response.is_none() && method == "GET" && url.starts_with("/collaboraUrl?server=") {
          let co_host = &url[21..url.len()].split(" ").next().unwrap_or("");

          let mut data = Vec::new();
          {
          let mut handle = Easy::new();
          _ = handle.url(format!("{}/hosting/discovery", co_host).as_str());
            let mut transfer = handle.transfer();
            _ = transfer.write_function(|new_data| {
              data.extend_from_slice(new_data);
              Ok(new_data.len())
            });
            _ = transfer.perform();
          }
          let xml_str = std::str::from_utf8(&data).unwrap_or("");
          if let Ok(package) = parser::parse(xml_str) {
            let xpath = "/wopi-discovery/net-zone/app[@name='text/plain']/action/@urlsrc";
            let document = package.as_document();
            if let Ok(online_url) = evaluate_xpath(&document, xpath) {
              response = Some(
                format!("{{\"url\":\"{}\",\"token\":\"test\"}}",
                online_url.string()));
            }
          }
        }
      
        if response.is_none() && method == "GET" && !url.contains("../") {
          let mut path : String = url.to_string();
          if path.ends_with("/") { path = format!("{path}/index.html") }
          response = Some(fs::read_to_string(format!("../html{path}")).unwrap_or_default());
        }
      
        if let Some(ref u_response) = response {
          let response = format!("HTTP/1.1 200 OK\r\n\r\n{u_response}");
          _ = stream.write_all(response.as_bytes());  // ignore err
        }
      }
    }
  }
}

const GET_FILE_JSON: &str = r#"{
  "BaseFileName": "test.txt",
  "Size": 11,
  "UserId": 1,
  "UserCanWrite": true,
  "EnableInsertRemoteImage": true
}"#;
