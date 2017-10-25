package p2p

import (
    "log"
    "net"
    "strconv"
)

func handleConnection(conn net.Conn) {
    buffer := make([]byte, 1024)

    for {
        _, err := conn.Read(buffer)
        if err != nil {
            log.Printf("Error reading from conn: %v\n", err)
            // conn.Write([]byte("bad"))
            continue
        }

        // TODO: Handle maintaining a peer cache here
        conn.Write([]byte("ok"))
    }
}

func listenerAcceptLoop(listener net.Listener) {
    defer listener.Close()

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Error accepting connection: %v\n", err)
            continue
        }

        log.Println("Accepting connection")
        go handleConnection(conn)
    }
}

func StartListener(port int) (error) {
    listener, err := net.Listen("tcp", "0.0.0.0:" + strconv.Itoa(port))
    if err != nil {
        return err
    }

    go listenerAcceptLoop(listener)
    return nil
}

