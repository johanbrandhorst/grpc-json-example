# grpc-json-example

Example of using gRPC-Go with JSON as the transport encoding.

## The JSON Codec

The included codec package contains [an implementation](./codec/json.go) of
[`encoding.Codec`](https://godoc.org/google.golang.org/grpc/encoding#Codec)
for JSON payloads.

## Server setup

The only thing necessary to make your gRPC server respond to gRPC requests
encoded with JSON is to import the relevant package (to let it register with
gRPC-Go's encoding registry):

```go
import _ "github.com/johanbrandhorst/grpc-json-example/codec"
```

## Request examples

### gRPC client

Using a gRPC Client, simply initiate using the correct content-subtype as a `grpc.DialOption`:

```go
grpc.WithDefaultCallOptions(grpc.CallContentSubtype(codec.JSON{}.Name()))
```

The [included client](./cmd/client/main.go) shows a full example.

### cURL

```bash
$ echo -en '\x00\x00\x00\x00\x17{"id":1,"role":"ADMIN"}' | curl -ss -k --http2 \
        -H "Content-Type: application/grpc+json" \
        -H "TE:trailers" \
        --data-binary @- \
        https://localhost:10000/example.UserService/AddUser | od -bc
0000000 000 000 000 000 002 173 175
         \0  \0  \0  \0 002   {   }
0000007
$ echo -en '\x00\x00\x00\x00\x17{"id":2,"role":"GUEST"}' | curl -ss -k --http2 \
        -H "Content-Type: application/grpc+json" \
        -H "TE:trailers" \
        --data-binary @- \
        https://localhost:10000/example.UserService/AddUser | od -bc
0000000 000 000 000 000 002 173 175
         \0  \0  \0  \0 002   {   }
0000007
$ echo -en '\x00\x00\x00\x00\x02{}' | curl -k --http2 \
        -H "Content-Type: application/grpc+json" \
        -H "TE:trailers" \
        --data-binary @- \
        --output - \
        https://localhost:10000/example.UserService/ListUsers
F{"id":1,"role":"ADMIN","create_date":"2018-07-21T20:18:21.961080119Z"}F{"id":2,"role":"GUEST","create_date":"2018-07-21T20:18:29.225624852Z"}
```

Explanation:

Using `cURL` to send requests requires manually adding the
[gRPC HTTP2 message payload header](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests)
to the payload:

```bash
'\x00\x00\x00\x00\x17{"id":1,"role":"ADMIN"}'
#<-->----------------------------------------- Compression boolean (1 byte)
#    <-------------->------------------------- Payload size (4 bytes)
#                    <--------------------->-- JSON payload
```

Headers must include `TE` and the correct `Content-Type`:
```bash
 -H "Content-Type: application/grpc+json" -H "TE:trailers"
```

The string after `application/grpc+` in the `Content-Type` header
must match the `Name()` of the codec registered in the server.

The endpoint must match the name of the name of the proto package,
the service and finally the method:

```bash
https://localhost:10000/example.UserService/AddUser
```

The responses are prefixed by the same header as the requests:

```bash
'\0  \0  \0  \0 002   {   }'
#<-->------------------------ Compression boolean (1 byte)
#    <------------>---------- Payload size (4 bytes)
#                    <---->-- JSON payload
```
