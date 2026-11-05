package network

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

// SendAndRecv sends a payload to the specified TCP server
// and returns the response received from it.
func SendAndRecv(address string, port string, payload string) ([]byte, error) {
	client, err := net.DialTimeout(
		"tcp",
		fmt.Sprintf("%s:%s", address, port),
		2*time.Second,
	)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	if err := client.SetWriteDeadline(time.Now().Add(30 * time.Second)); err != nil {
		return nil, err
	}
	_, err = client.Write(append([]byte(payload), '\n'))
	if err != nil {
		return nil, err
	}
	if err := client.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
		return nil, err
	}
	reader := bufio.NewReader(client)
	response, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return response, nil
}
