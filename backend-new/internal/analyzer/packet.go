package analyzer

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type FlowKey struct {
	SrcIP, DstIP     string
	SrcPort, DstPort uint16
	Protocol         string
}

type FlowStats struct {
	FlowKey
	Packets  uint64
	Bytes    uint64
	LastSeen time.Time
}

type ConnTracker struct {
	iface    string
	handle   *pcap.Handle
	flows    map[FlowKey]*FlowStats
	mu       sync.RWMutex
	stopChan chan struct{}
	stopped  chan struct{}
}

func NewConnTracker(iface string) *ConnTracker {
	return &ConnTracker{
		iface:    iface,
		flows:    make(map[FlowKey]*FlowStats),
		stopChan: make(chan struct{}),
		stopped:  make(chan struct{}),
	}
}

func (ct *ConnTracker) Start() error {
	handle, err := pcap.OpenLive(ct.iface, 1600, true, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("pcap open: %w", err)
	}
	ct.handle = handle

	go ct.capture()
	return nil
}

func (ct *ConnTracker) Stop() {
	close(ct.stopChan)
	<-ct.stopped
	ct.handle.Close()
}

func (ct *ConnTracker) Stats() []FlowStats {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	snapshot := make([]FlowStats, 0, len(ct.flows))
	for _, f := range ct.flows {
		snapshot = append(snapshot, *f)
	}
	return snapshot
}

func (ct *ConnTracker) capture() {
	packetSource := gopacket.NewPacketSource(ct.handle, ct.handle.LinkType())
	for {
		select {
		case <-ct.stopChan:
			close(ct.stopped)
			return
		case pkt := <-packetSource.Packets():
			ct.process(pkt)
		}
	}
}

func (ct *ConnTracker) process(pkt gopacket.Packet) {
	network := pkt.NetworkLayer()
	transport := pkt.TransportLayer()
	if network == nil || transport == nil {
		return
	}

	src, dst := network.NetworkFlow().Endpoints()
	srcPort, dstPort := transport.TransportFlow().Endpoints()
	key := FlowKey{
		SrcIP:    src.String(),
		DstIP:    dst.String(),
		SrcPort:  parsePort(srcPort.String()),
		DstPort:  parsePort(dstPort.String()),
		Protocol: transport.LayerType().String(),
	}

	ct.mu.Lock()
	defer ct.mu.Unlock()

	stats, exists := ct.flows[key]
	if !exists {
		stats = &FlowStats{FlowKey: key}
		ct.flows[key] = stats
	}
	stats.Packets++
	stats.Bytes += uint64(len(pkt.Data()))
	stats.LastSeen = time.Now()
}

func parsePort(p string) uint16 {
	port, _ := strconv.Atoi(p)
	return uint16(port)
}
