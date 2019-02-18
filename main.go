package main

import (
    "fmt" // For printing out to terminal
    "net" // For network functionality (opening sockets)
    "bufio" // For reading in data from sockets
    "strings" // For manipulating received data
    "io/ioutil" // For dealing with files
)

// The path to start looking for files from.
const PATH = "./data"
const INDEX = "/index.txt"

/**
 * The main entry point of our application.
 */
func main() {
    // Bind a TCP/IP socket to port 8080
    socket, err := net.Listen("tcp", ":8080")

    if err != nil { return }

    // Report to user that we're ready to connect.
    fmt.Println("Ready to accept connections!")

    for {
        // Accept a connection
        connection, err := socket.Accept()

        if err != nil { break }

        // Report that the connection was accepted.
        fmt.Println("Connection accepted!")

        // Handle the connection we just received by spinning up a goroutine.
        go HandleConnection(connection)
    }
}

/**
 * Handles a TCP/IP connection made to the server.
 */
func HandleConnection(connection net.Conn) {
    // Create a scanner from the connection we received.
    scanner := bufio.NewScanner(connection)

    // Read in the request header.
    no_err := scanner.Scan()

    // If we can't read in the header, abort and notify the client that we
    // received a bad request.
    if !no_err {
        fmt.Fprintln(connection, "HTTP/1.0 400 Bad Request\r\n\r\n")
        connection.Close()
        return
    }

    // Retreive the request header from our scanner.
    request_head := scanner.Text()

    // Read in the next lines from the request. If there is no err, report the
    // retreived line of the request. If there is an error, abort. Stop
    // reading if we have reached the end of the request (we get an empty
    // string from scanner.Text()).
    for {
        no_err := scanner.Scan()

        if !no_err { break }

        line := scanner.Text()

        fmt.Println(line)

        if line == "" { break }
    }

    // Get the request type, request location, and HTTP version from the
    // request header.
    tokens := strings.Split(request_head, " ")

    request_type := tokens[0]
    // location     := tokens[1]
    // version      := tokens[2]

    switch request_type {
    case "GET":
        {
            fmt.Println("Get request received.")

            data, err := ioutil.ReadFile(fmt.Sprintf("%s%s", PATH, INDEX))

            if err != nil {
                fmt.Fprintln(connection, "HTTP/1.0 404 Not Found\r\n\r\n")
                connection.Close()
                return
            }

            fmt.Fprintln(connection,
                         fmt.Sprintf("HTTP/1.0 200 OK\r\n\r\n%s", string(data)))

            break
        }
    }

    // Declare that there is no more data to be read.
    fmt.Println("End of Connection.")

    // Close the connection now that we've done what we wanted to.
    connection.Close()
}
