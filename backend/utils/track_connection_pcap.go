package utils

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type ConnectionInfo struct {
	Protocol    string    `json:"protocol"`
	SourceIP    string    `json:"source_ip"`
	DestIP      string    `json:"dest_ip"`
	SourcePort  string    `json:"source_port"`
	DestPort    string    `json:"dest_port"`
	PacketCount string    `json:"packet_count"`
	ByteCount   string    `json:"byte_count"`
	LastSeen    time.Time `json:"last_seen"`
}

type ConnectionTracker struct {
	connections map[string]*ConnectionInfo
	mutex       sync.RWMutex
	handle      *pcap.Handle
	cancel      context.CancelFunc
	ctx         context.Context
	isRunning   bool
}

func NewConnectionTracker() *ConnectionTracker {
	ctx, cancel := context.WithCancel(context.Background())
	return &ConnectionTracker{
		connections: make(map[string]*ConnectionInfo),
		ctx:         ctx,
		cancel:      cancel,
		isRunning:   false,
	}
}

func (ct *ConnectionTracker) updateConnection(key string, conn *ConnectionInfo) {
	ct.mutex.Lock()
	defer ct.mutex.Unlock()

	if existing, exists := ct.connections[key]; exists {
		packetCount := parseInt(existing.PacketCount) + 1
		byteCount := parseInt(existing.ByteCount) + parseInt(conn.ByteCount)
		existing.PacketCount = fmt.Sprintf("%d", packetCount)
		existing.ByteCount = fmt.Sprintf("%d", byteCount)
		existing.LastSeen = time.Now()
	} else {
		conn.LastSeen = time.Now()
		ct.connections[key] = conn
	}
}

func (ct *ConnectionTracker) GetConnections() []ConnectionInfo {
	ct.mutex.RLock()
	defer ct.mutex.RUnlock()

	var result []ConnectionInfo
	now := time.Now()
	for key, conn := range ct.connections {
		// 5분 이상 업데이트 되지 않은 연결은 제거
		if now.Sub(conn.LastSeen) > 5*time.Minute {
			delete(ct.connections, key)
			continue
		}
		result = append(result, *conn)
	}
	return result
}

func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}

// StartCapture starts packet capture on the specified interface
func (ct *ConnectionTracker) StartCapture(interfaceName string) error {
	if ct.isRunning {
		return fmt.Errorf("packet capture is already running")
	}

	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("error opening interface %s: %w", interfaceName, err)
	}

	ct.handle = handle
	ct.isRunning = true

	go ct.capturePackets()

	return nil
}

// StopCapture stops the current packet capture
func (ct *ConnectionTracker) StopCapture() error {
	if !ct.isRunning {
		return fmt.Errorf("packet capture is not running")
	}

	ct.cancel()
	if ct.handle != nil {
		ct.handle.Close()
	}

	ct.isRunning = false
	return nil
}

// IsRunning returns the current capture status
func (ct *ConnectionTracker) IsRunning() bool {
	return ct.isRunning
}

func (ct *ConnectionTracker) capturePackets() {
	defer ct.handle.Close()

	packetSource := gopacket.NewPacketSource(ct.handle, ct.handle.LinkType())

	for {
		select {
		case <-ct.ctx.Done():
			return
		case packet := <-packetSource.Packets():
			if packet == nil {
				continue
			}

			ipLayer := packet.Layer(layers.LayerTypeIPv4)
			if ipLayer == nil {
				continue
			}

			ip, _ := ipLayer.(*layers.IPv4)

			var protocol string
			var srcPort, dstPort string

			tcpLayer := packet.Layer(layers.LayerTypeTCP)
			if tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)
				protocol = "tcp"
				srcPort = tcp.SrcPort.String()
				dstPort = tcp.DstPort.String()
			}

			udpLayer := packet.Layer(layers.LayerTypeUDP)
			if udpLayer != nil {
				udp, _ := udpLayer.(*layers.UDP)
				protocol = "udp"
				srcPort = udp.SrcPort.String()
				dstPort = udp.DstPort.String()
			}

			if protocol == "" {
				continue
			}

			connInfo := &ConnectionInfo{
				Protocol:    protocol,
				SourceIP:    ip.SrcIP.String(),
				DestIP:      ip.DstIP.String(),
				SourcePort:  srcPort,
				DestPort:    dstPort,
				PacketCount: "1",
				ByteCount:   fmt.Sprintf("%d", len(packet.Data())),
			}

			key := fmt.Sprintf("%s-%s:%s-%s:%s",
				protocol, ip.SrcIP, srcPort, ip.DstIP, dstPort)

			ct.updateConnection(key, connInfo)
		}
	}
}

// 전역 ConnectionTracker 인스턴스
var globalTracker *ConnectionTracker
var once sync.Once

// InitializeCapture initializes and starts packet capture
func InitializeCapture(interfaceName string) error {
	var err error
	once.Do(func() {
		globalTracker = NewConnectionTracker()
		err = globalTracker.StartCapture(interfaceName)
	})
	return err
}

// StopCapture stops the global packet capture
func StopCapture() error {
	if globalTracker == nil {
		return fmt.Errorf("packet capture has not been initialized")
	}
	return globalTracker.StopCapture()
}

// GetCurrentConnections returns current connection information
func GetCurrentConnections() ([]ConnectionInfo, error) {
	if globalTracker == nil {
		return nil, fmt.Errorf("packet capture has not been initialized")
	}
	return globalTracker.GetConnections(), nil
}
