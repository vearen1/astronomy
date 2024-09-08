package server

import (
	"astronomy/astronomy/internal/packet"
	"astronomy/astronomy/internal/protocol"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

type MinecraftServer struct {
	Address string
	Port    int
	Status  ServerStatus
}

type ServerStatus struct {
	Version     ServerVersion `json:"version"`
	Players     PlayerInfo    `json:"players"`
	Description string        `json:"description"`
}

type ServerVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type PlayerInfo struct {
	Max    int `json:"max"`
	Online int `json:"online"`
}

func NewMinecraftServer(address string, port int) *MinecraftServer {
	return &MinecraftServer{
		Address: address,
		Port:    port,
		Status: ServerStatus{
			Version: ServerVersion{
				Name:     "1.21",
				Protocol: 767,
			},
			Players: PlayerInfo{
				Max:    100,
				Online: 0,
			},
			Description: "§bWelcome to our Minecraft server!",
		},
	}
}

func (s *MinecraftServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Address, s.Port))
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("Minecraft server listening on %s:%d", s.Address, s.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *MinecraftServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		_, err := protocol.ReadVarInt(conn)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading packet length: %v", err)
			}
			return
		}

		packetID, err := protocol.ReadVarInt(conn)
		if err != nil {
			log.Printf("Error reading packet ID: %v", err)
			return
		}

		switch packet.PacketID(packetID) {
		case packet.HandshakePacketID:
			packet := &packet.HandshakePacket{}
			if err := packet.Decode(conn); err != nil {
				log.Printf("Error decoding handshake packet: %v", err)
				return
			}
			if packet.NextState == 1 {
				s.handleStatusRequest(conn)
			} else {
				log.Printf("Unsupported next state: %d", packet.NextState)
				return
			}
		default:
			log.Printf("Unsupported packet ID: %d", packetID)
			return
		}
	}
}

func (s *MinecraftServer) handleStatusRequest(conn net.Conn) {
	// Read status request packet
	_, err := protocol.ReadVarInt(conn)
	if err != nil {
		log.Printf("Error reading status request packet length: %v", err)
		return
	}

	packetID, err := protocol.ReadVarInt(conn)
	if err != nil {
		log.Printf("Error reading status request packet ID: %v", err)
		return
	}

	if packet.PacketID(packetID) != packet.StatusRequestPacketID {
		log.Printf("Expected status request packet, got: %d", packetID)
		return
	}

	// Send status response
	statusJSON, err := json.Marshal(s.Status)
	if err != nil {
		log.Printf("Error marshaling status JSON: %v", err)
		return
	}

	responsePacket := &packet.StatusResponsePacket{Response: string(statusJSON)}
	responseData, err := responsePacket.Encode()
	if err != nil {
		log.Printf("Error encoding status response packet: %v", err)
		return
	}

	if err := protocol.WriteVarInt(conn, int32(len(responseData))); err != nil {
		log.Printf("Error writing response packet length: %v", err)
		return
	}

	if _, err := conn.Write(responseData); err != nil {
		log.Printf("Error writing response packet: %v", err)
		return
	}

	// Handle ping request
	s.handlePingRequest(conn)
}

func (s *MinecraftServer) handlePingRequest(conn net.Conn) {
	_, err := protocol.ReadVarInt(conn)
	if err != nil {
		log.Printf("Error reading ping request packet length: %v", err)
		return
	}

	packetID, err := protocol.ReadVarInt(conn)
	if err != nil {
		log.Printf("Error reading ping request packet ID: %v", err)
		return
	}

	if packet.PacketID(packetID) != packet.PingRequestPacketID {
		log.Printf("Expected ping request packet, got: %d", packetID)
		return
	}

	pingPacket := &packet.PingRequestPacket{}
	if err := pingPacket.Decode(conn); err != nil {
		log.Printf("Error decoding ping request packet: %v", err)
		return
	}

	// Send ping response
	responsePacket := &packet.PingResponsePacket{Payload: pingPacket.Payload}
	responseData, err := responsePacket.Encode()
	if err != nil {
		log.Printf("Error encoding ping response packet: %v", err)
		return
	}

	if err := protocol.WriteVarInt(conn, int32(len(responseData))); err != nil {
		log.Printf("Error writing ping response packet length: %v", err)
		return
	}

	if _, err := conn.Write(responseData); err != nil {
		log.Printf("Error writing ping response packet: %v", err)
		return
	}
}
