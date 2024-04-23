In order to create multiple instances of the server, 
we can make it get the ip address and port from the
console arguments.
```
if len(os.Args) < 3 {
    fmt.Println("Usage: go run server.go <ip_address> <port>")
    return
}

ipAddress := os.Args[1]
port := os.Args[2]
addr := fmt.Sprintf("%s:%s", ipAddress, port)

log.Println("Starting HTTP server on port 8080...")
if err := fasthttp.ListenAndServe(addr, requestHandler); err != nil {
    log.Fatalf("Error starting HTTP server: %v", err)
}
```

Now the server can run on different ip addresses and the same port:
```go run main.go 127.0.0.2 8080```

```go run main.go 127.0.0.3 8080```

In order to map a domain name to one or multiple ip addresses, 
we need to edit the hosts file.

On Windows, the file can be found at `C:\Windows\System32\drivers\etc\hosts`.

We add our ip addresses and map them to the same domain name as follows:
```
127.0.0.2       adtelligent-internship.com
127.0.0.3       adtelligent-internship.com
```

Then we input the following command to flush the DNS cache: `ipconfig /flushdns`

Now instead of the ip address, we can ust the domain name.

```
curl "http://adtelligent-internship.com:8080/campaigns_slice?source_id=2&domain=gmail.com"


StatusCode        : 200
StatusDescription : OK
Content           : [{"ID":27,"Name":"Campaign
                    27","FilterType":"black","Domains":["aol.com","mail.ru","orange.fr","yahoo.com"],"SourceID":2}]
RawContent        : HTTP/1.1 200 OK
                    Content-Length: 122
                    Content-Type: application/json
                    Date: Tue, 23 Apr 2024 15:46:08 GMT
                    Server: fasthttp

                    [{"ID":27,"Name":"Campaign 27","FilterType":"black","Domains":["aol.com",...
Forms             : {}
Headers           : {[Content-Length, 122], [Content-Type, application/json], [Date, Tue, 23 Apr 2024 15:46:08 GMT],
                    [Server, fasthttp]}
Images            : {}
InputFields       : {}
Links             : {}
ParsedHtml        : mshtml.HTMLDocumentClass
RawContentLength  : 122
```

I noticed that the requests are sent only to the first instance. 
I guess it is due to the hosts file priority. 

When I shut the first instance down, the requests are sent to the second instance.
