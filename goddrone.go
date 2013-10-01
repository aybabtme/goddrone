package goddrone

import (
	"fmt"
	"net"
)

const (
	droneIP = "192.168.1.1"
	udpNet  = "udp"
	udpPort = "5559"
)

type Drone struct {
	addr *net.UDPAddr
	conn *net.UDPConn
}

func Connect(ip string) (*Drone, error) {

	droneAddr, err := net.ResolveUDPAddr(udpNet, ip+":"+udpPort)
	if err != nil {
		return nil, fmt.Errorf("Failed resolving drone address, %v", err)
	}

	conn, err := net.DialUDP(udpNet, nil, droneAddr)
	if err != nil {
		return nil, fmt.Errorf("Dial failed, %v", err)
	}

	return &Drone{droneAddr, conn}, nil
}

func (d *Drone) send(b []byte) error {
	n, err := d.conn.Write(b)
	if err != nil {
		return err
	}
	if n != len(b) {
		return fmt.Errorf("want %d bytes, sent %d", len(b), n)
	}
	return nil
}

func (d *Drone) Disconnect() error {
	return d.conn.Close()
}
