package broadcast

import (
	"errors"
	"fmt"
	"net"
	"time"
)

var (
	broadcastAddress = "255.255.255.255"
	network          = "udp4"
)

//UDPBroadcaster is for broadcasting on a given
type UDPBroadcaster struct {
	port     int
	timeout  int
	payload  string
	findself bool
}

//NewUDPBroadcaster returns an intialized broadcast object for the given port and payload.
func NewUDPBroadcaster(port int, payload string) *UDPBroadcaster {
	return &UDPBroadcaster{
		port:    port,
		payload: payload,
		timeout: 10,
	}
}

//SetFindself controls if this library should find it's own listener in the results of a Discover().
func (bc *UDPBroadcaster) SetFindself(setting bool) {
	bc.findself = setting
}

//SetTimeout sets the timeout for how long to wait for an answer
func (bc *UDPBroadcaster) SetTimeout(setting int) {
	bc.timeout = setting
}

//Discover listens for responses to a UDP Broadcast sent to the LAN broadcast address, and returns any responders IP addresses.
func (bc *UDPBroadcaster) Discover() (ips []string, err error) {

	// Check defaults
	if bc.port == 0 {
		return nil, errors.New("invalid port to discover")
	}
	if bc.payload == "" {
		return nil, errors.New("no payload supplied")
	}

	// Variables for later
	var serverAddr *net.UDPAddr
	var ourAddr *net.UDPAddr
	var conn *net.UDPConn
	inBuf := make([]byte, 1024)
	// get server address
	server := fmt.Sprintf("%s:%d", broadcastAddress, bc.port)
	serverAddr, err = net.ResolveUDPAddr(network, server)
	if err != nil {
		return
	}
	// get local address
	ourAddr, err = net.ResolveUDPAddr(network, fmt.Sprintf(":%d", bc.port))
	if err != nil {
		return
	}
	// get local connection
	conn, err = net.ListenUDP(network, ourAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	// send the Discover message
	discoverMsg := []byte(bc.payload)
	if _, err = conn.WriteTo(discoverMsg, serverAddr); err != nil {
		return
	}

	// Get the localIPs on this system as we want to parse out only non-local responses.
	locals := getLocalIPs()

	// Try to read as many as possible.
	addresses := make(chan *net.UDPAddr)
	errs := make(chan error)
	cancel := make(chan bool, 1)
	go func() {
		for {
			// see if we have been canceled...
			select {
			case <-cancel:
				return
			default:
			}
			var err error
			var fromAddr *net.UDPAddr
			_, fromAddr, err = conn.ReadFromUDP(inBuf)
			if err != nil {
				errs <- err
				return
			}
			if bc.findself {
				// Add all ips IPs to the ips results
				addresses <- fromAddr
			} else {
				// see if we have a match to a local host
				isours := false
				for ip := range locals {
					if ip == fromAddr.IP.String() {
						isours = true
					}
				}
				if !isours {
					addresses <- fromAddr
				}
			}

		}
	}()

	// Drain the channels
	for {
		select {
		case a := <-addresses:
			ips = append(ips, a.IP.String())
		case err = <-errs:
			cancel <- true
			return nil, err
		case <-time.After(2 * time.Second):
			cancel <- true
			return
		}
	}

	return
}
