package goddrone

import (
	"net"
	"testing"
)

func TestCanConnect(t *testing.T) {

	laddr, err := net.ResolveUDPAddr(udpNet, ":"+udpPort)
	if err != nil {
		t.Fatalf("Resolve, %v", err)
	}

	ln, err := net.ListenUDP(udpNet, laddr)
	if err != nil {
		t.Fatalf("Listen, %v", err)
	}
	defer ln.Close()

	want := "lol."

	recvC := make(chan string, 1)
	errC := make(chan error)

	go func() {
		b := make([]byte, len(want))
		n, _, err := ln.ReadFromUDP(b)

		recvC <- string(b[:n])

		if err != nil {
			errC <- err
		}
		close(recvC)
	}()

	go func() {
		drone, err := Connect("127.0.0.1")
		if err != nil {
			t.Fatalf("Connect, %v", err)
		}
		defer drone.Disconnect()

		err = drone.send([]byte(want))
		if err != nil {
			errC <- err
		}
	}()

	got := <-recvC

	close(errC)

	for err := range errC {
		t.Errorf("goroutine reported error, %v", err)
	}

	if got != want {
		t.Errorf("Want %v, got %v", want, got)
	}

}
