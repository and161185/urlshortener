# URL Shortener

## Description

**URL Shortener** is a web service written in Go that allows users to shorten long URLs. Users can enter long links and receive short URLs that redirect to the original addresses.

## Features

- **URL Shortening**: Enter a long URL and get a short link.
- **Redirection**: Accessing a short link redirects to the original URL.
- **Statistics**: Track the number of visits to each shortened link (if implemented).

## Installation and Running

1. **Clone the repository**:

   ```bash
   git clone https://github.com/and161185/urlshortener.git
   cd urlshortener
   ```

2. **Install dependencies**:

   Ensure that you have Go version 1.16 or higher installed.

   ```bash
   go mod download
   ```

3. **Run the application**:

   ```bash
   go run main.go
   ```

   By default, the server runs on port `8080`. You can change the port by setting the `PORT` environment variable.

## Usage

- **Shorten URL**:

  Send a POST request to `/shorten` with a JSON body containing the `url` field:

  ```json
  {
    "url": "https://example.com/very/long/url"
  }
  ```

  The response will contain a JSON object with the shortened link:

  ```json
  {
    "short_url": "http://localhost:8080/abc123"
  }
  ```

- **Redirection**:

  Visit the shortened link, and you will be redirected to the original URL.

## Configuration

- **Port**: By default, the application listens on port `8080`. To change the port, set the `PORT` environment variable:

  ```bash
  export PORT=9090
  ```

## Testing

To run tests, execute:

```bash
go test ./...
```

## Dependencies

- [gorilla/mux](https://github.com/gorilla/mux): HTTP request routing.
- [go-redis/redis](https://github.com/go-redis/redis): Redis client (if used for storage).

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

