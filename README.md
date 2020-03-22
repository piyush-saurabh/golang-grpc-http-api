# Golang gRPC/HTTP API

This repository contains a microservice developed using golang and communicates via gRPC. The microservice is exposed as REST API. 

Ref: https://medium.com/@amsokol.com/tutorial-how-to-develop-go-grpc-microservice-with-http-rest-endpoint-middleware-kubernetes-daebb36a97e9

Package used for exposing gRPC service as REST API: https://github.com/grpc-ecosystem/grpc-gateway

![Markdown Logo](https://github.com/piyush-saurabh/golang-grpc-http-api/blob/master/grpc-http-api.png)

## Project Structure
Ref: https://github.com/golang-standards/project-layout

| Package Name  | Purpose             |
|---------------|-------------------|
|/api |All the **.proto** files will be here. Protocol definition files, OpenAPI/Swagger specs, JSON schema files.
|/pkg |All the **compiled** proto files will be here. Put the re-usable code here.
|/cmd |It has **main** functions for both client and server.
|/third-party | 3rd party libraries

## Code Snippets

### Modify proto file to expose gRPC service as REST API

```proto
import "google/api/annotations.proto";
import "protoc-gen-swagger/options/annotations.proto";

option (grpc.gateway.protoc_gen_swagger.options.openapiv2_swagger) = {
	info: {
		title: "ToDo service";
		version: "1.0";
		contact: {
			name: "golang-grpc-http-api project";
			url: "https://github.com/piyush-saurabh/golang-grpc-http-api";
			email: "test@email.com";
        };
    };
    schemes: HTTP;
    consumes: "application/json";
    produces: "application/json";
    responses: {
		key: "404";
		value: {
			description: "Returned when the resource does not exist.";
			schema: {
				json_schema: {
					type: STRING;
				}
			}
		}
	}
};

// Service to manage list of todo tasks
service ToDoService {
    // This gRPC service can be accessed via REST client at path /v1/to/all
    rpc ReadAll(ReadAllRequest) returns (ReadAllResponse){
        option (google.api.http) = {
            get: "/v1/todo/all"
        };
}
```

### Start the REST Gateway
Rest gateway is **HTTP Listener + gRPC Client**

```go
import "github.com/grpc-ecosystem/grpc-gateway/runtime"

func RunServer(ctx context.Context, grpcPort, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create a HTTP request router/multiplexer to handle route based on URL path
	mux := runtime.NewServeMux()

  // Create the gRPC client to make call to gRPC server
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// On receiving HTTP request, make a HTTP/2 call to gRPC server running on the localhost (if they are the part of same binary)
	// RegisterToDoServiceHandlerFromEndpoint is the HTTP/REST handler created automatically grpc-gateway
	if err := v1.RegisterToDoServiceHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts); err != nil {
		log.Fatalf("failed to register gRPC handler: %v", err)
	}

	// Serve the HTTP endpoint
	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
		}
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()
  // Start the HTTP Listener
	return srv.ListenAndServe()
}
```

### Start gRPC and HTTP server as the part of same binary
This code is part of server's main()
```go
ctx := context.Background()

// Start HTTP server as a goroutine
go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
}()
server := grpc.NewServer()
v1.RegisterToDoServiceServer(server, v1API)

// start gRPC server
server.Serve(listen)
  
```

