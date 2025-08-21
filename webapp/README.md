# Content

- [dotNET](dotNET): .NET example
- [golang](golang): Go example
- [nodejs](nodejs): NodeJS example
- [nodejs-esm](nodejs-esm): same as the nodejs example but using ESM
- [php](php): PHP example
- [python](python): Python example
- [reactjs](reactjs): ReactJS example
- [rust](rust): Rust example
- html: The frontend part used by some other examples. This is NOT an example itself.

# Running the WebApp Examples

These examples demonstrate how to [integrate Collabora Online](https://sdk.collaboraonline.com/docs/How_to_integrate.html) via the [WOPI API](https://en.wikipedia.org/wiki/Web_Application_Open_Platform_Interface) into other applications. See the example folders for specific instructions.

Serve via HTTPS/TLS/SSL to make web browsers allow [secure context](https://developer.mozilla.org/docs/Web/Security/Secure_Contexts) actions like using the clipboard. Set the environment variables `SSL_CRT_FILE=crt.pem` and `SSL_KEY_FILE=key.pem` to enable HTTPS when running the examples.

The example will run (in) an HTTP(S) server. To then make use of it, you additionally need a Collabora Online HTTP(S) Server. And the example server and the Collabora Online Server need to be able to reach each other to work. The Collabora Online Server will contact the example server by the same address you use to open the example's web page.

See also: https://sdk.collaboraonline.com/docs/examples.html

## Using a public Collabora Online Server

When running an example, you might enter `https://demo.eu.collaboraonline.com` or `https://demo.us.collaboraonline.com` as Collabora Online Server in the example's web form.  
This requires serving the example on a public address with a valid certificate for HTTPS. **Self-Signed certificates won't work!**

E.g. serve the example publicly at `https://yourdomain.example.org/`. And (**important**) use `https://yourdomain.example.org/` to open the example in your web browser. Do NOT open the example using `https://127.0.0.1/` or `https://SOMETHING_INTERNAL/`, because the address from your browsers URL bar will be used by the Collabora Online Server to to callbacks on the example.

You might either run the example directly via HTTPS with a valid certificate using the environment variables `SSL_CRT_FILE` and `SSL_KEY_FILE`. Or you run it with plain HTTP and put it behind a [reverse proxy](https://en.wikipedia.org/wiki/Reverse_proxy) like [Apache](https://httpd.apache.org/docs/current/howto/reverse_proxy.html) to enable HTTPS.

Do NOT put the example behind any kind of authentication like HTTP Basic Auth, because the Collabora Online Server need to do callbacks on the example server. For additional security the `demo.*.collaboraonline.com` servers sign their data using [WOPI proof](https://sdk.collaboraonline.com/docs/advanced_integration.html#wopi-proof), so your code might verify these signatures.

## Run it locally

To have a secure context (see above), both the example and the Collabora Online Server must have HTTPS enabled.

### Run the example locally

To run the example with HTTPS you might create a self-signed X.509 certificate using [minica](https://github.com/jsha/minica) or running:  
`openssl req -x509 -newkey rsa:4096 -keyout key.pem -out crt.pem -days 365 -nodes -subj "/CN=127.0.0.1"`  

**DO NOT DISABLE CERTIFICATE VERIFICATION IN PRODUCTION!**  
To make the example server accept a self-signed certificate used by the Collabora Online Server, set the environment variable `DISABLE_TLS_CERT_VALIDATION=1`.

### Build Collabora Online from source

You might [compile the Collabora Online Server from source](https://collaboraonline.github.io/post/build-code/) and run it locally. When doing this, you won't need to make the example server publicly reachable.

To use HTTPS do NOT disable SSL when building Collabora Online.

**DO NOT DISABLE CERTIFICATE VERIFICATION IN PRODUCTION!**  
After building open the file `coolwsd.xml` and ensure `config -> ssl -> ssl_verification` is set to `false` to accept a self-signed certificate used by the example server.

### Download pre-build Collabora Online Development Edition

see: https://www.collaboraonline.com/code/#learnmorecode

SSL options for the container (Docker) image:
```
--o:ssl.key_file_path="key.pem" \
--o:ssl.cert_file_path="crt.pem" \
--o:ssl.ca_file_path="ca-chain.cert.pem" 
```

## Miscellaneous

Most examples will accept these environment variables:
- `SSL_CRT_FILE` & `SSL_KEY_FILE`: Serve with HTTPS.
  - Without HTTP will be used and things like clipboard support will be disabled by your browser.
- `DISABLE_TLS_CERT_VALIDATION=1`: Don't validate the Collabora Server's HTTPS certificate. NOT FOR PRODUCTION!
- `PORT`: Most examples use TCP port `3000` as default.

# Setup Schematic

The example implements the frontend as well as the document server backend. In a real application, the frontend and document server backend might be different servers.

Simplified sequence diagram:

```
                           example              Collabora Online
web browser             HTTP(S) server           HTTP(S) Server
-----------             --------------          ----------------
     |                        |                         |
     |  calls: frontend_addr  |                         |
     |  data: collabora_addr  |                         |
     | ---------------------> |                         |
     |                        |                         |
     |              ----------------------              |
     |              | doc_srv_addr=      |              |
     |              | frontend_addr      |              |
     |              | # addr of this srv |              |
     |              ----------------------              |
     |                        |                         |
     |                        |  calls: collabora_addr  |
     |                        |  data: doc_srv_addr     |
     |                        | ----------------------> |
     |                        |                         |
                ...                      ...
     |                        |                         |
     |                        |   calls: doc_srv_addr   |
     |                        | <---------------------- |
```
