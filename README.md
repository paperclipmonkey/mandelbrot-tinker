# mandelbrot-tinker
Toy challenge to learn Go.

![Mandelbrot Set](points.png)

## Build Status

![Go Build](https://github.com/paperclipmonkey/mandelbrot-tinker/actions/workflows/go.yml/badge.svg)
![Docker Build](https://github.com/paperclipmonkey/mandelbrot-tinker/actions/workflows/build.yml/badge.svg)

## How to Use

### GCP
See it running on GCP at:
[https://mandelbrot-629070084289.europe-west2.run.app/](https://mandelbrot-629070084289.europe-west2.run.app/)

### Run from a Docker Container
To run the application from a Docker container, use the following command:
```sh
docker run --pull=always --rm 59vkckvlkjdfglkjdfv/mandelbrot-tinker:latest
```
then open [http://localhost:80](http://localhost:80) in your web browser to interact with the Mandelbrot image.

### Build
To build the application, you need the following tools installed:
- Go 1.23 or later
- Docker

#### Required Tools
- [Go](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)

#### How to Develop
1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/mandelbrot-tinker.git
    cd mandelbrot-tinker
    ```
2. Install dependencies:
    ```sh
    go get
    ```
3. Run tests:
    ```sh
    go test ./...
    ```
4. Build the application:
    ```sh
    go build -o mandelbrot-tinker
    ```
    or
    ```
    go run .
    ```
