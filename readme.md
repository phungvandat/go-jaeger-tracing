# EXAMPLE JAEGER TRACING

![jaeger-go](https://raw.githubusercontent.com/phungvandat/jaegertracing/master/images/jaeger.png)

## SETUP

- `make setup`

## HOW TO USE

- Start server 1 (HTTP1, GRPC)

```
make server1
```

- Start server 2 (HTTP1)

```
make server2
```

- Make request from client to server1 to trace

```
make client
```

- Can call any http request to server1 or server2
- Frontend: `http://localhost:16686`
