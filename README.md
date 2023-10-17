# HTTP Server with REST API for Binary File CRUD Operations

This is a simple HTTP server implementation in Go, designed to provide a REST API for performing
CRUD (Create, Read, Update, Delete) operations on a single resource bina.
The resource data is stored in a binary file.

## API Endpoints

The server provides the following REST endpoints:

### GET /readyz

Health check for service, retrieve OK.

+ Return Http status code 200

### GET /records/{id:[0-9]+}

Retrieve a specific record from binary file by ID.

+ Return Http status code 200
+ Record formatted as JSON

### POST /records

Create a new record formatted as JSON.

+ Content-Type: application/json
+ Return Http status code 201

```
{
  "IntValue": 42,
  "StrValue": "foo",
  "BoolValue": true,
  "TimeValue": "2023-10-10T21:57:00+02:00"
}
```

### PUT /records/{id:[0-9]+}

Update an existing record by ID.

+ Content-Type: application/json
+ Return Http status code 200

```
{
  "IntValue": 42,
  "StrValue": "foo",
  "BoolValue": true,
  "TimeValue": "2023-10-10T21:57:00+02:00"
}
```

### DELETE /records/{id:[0-9]+}

Delete a record by ID (Only set ID to 0, other data are not changed)

+ Return Http status code 204

## Running the Server

The server will be accessible at http://localhost:8080.

```
go run main.go
```

Optional environment variable

+ PORT - specific server port (default value: 8080)
+ LOG_DEBUG - set log level to debug (default value: false)
+ BINARY_FILE_PATH - set path for binary file storage (default value: ./records.bin)

## Testing

You can run the unit tests using the following command:

```
go test ./...
```

## Docker Support

This project includes Docker support for containerization. You can build a Docker image and run the server within a container using the provided Dockerfile. Here are the steps to build and run the Docker container:

Build the Docker image:

```
docker build -t interview-test .
```

Run the Docker container:

```
docker run -p 8080:8080 interview-test
```

Run the Docker container with custom environment variable:

```
docker run --env-file ./docker_env -p 8080:8080 interview-test
```

The server will be accessible at http://localhost:8080 within the Docker container.

