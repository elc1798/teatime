package p2p

import (
    "net"
)

/*
 * Send specified data to the specified connection
 *
 * Returns the number of bytes sent, and an error if unsuccessful
 */
func SendData(conn *net.TCPConn, bytes []byte) (int, error) {
    return 0, nil
}

/*
 * Reads data from the specified connection
 *
 * Returns a byte array containing the read data,  number of bytes read and an
 * error if unsuccessful
 */
func ReadData(conn *net.TCPConn) ([]byte, int, error) {
    return make([]byte, 0), 0, nil
}

