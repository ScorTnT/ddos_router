package protect_manager

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ScorTnT/ddos_router/backend/config"
)

type AlertScanner interface {
	StartScan()
	StopScan()
	GetChannel() <-chan string
}

type ProtectManager struct {
	ttl           time.Duration
	blocksMu      sync.RWMutex
	blocks        map[string]time.Time
	addChannel    chan string
	deleteChannel chan string
	refreshTick   time.Duration
	alertScanner  AlertScanner
	ctx           context.Context
	cancel        context.CancelFunc
}

type BlockedIP struct {
	IP          string        `json:"ip"`
	ExpireAt    time.Time     `json:"expire_at"`
	OriginalTTL time.Duration `json:"original_ttl"`
}

func NewProtectManager(ttl time.Duration, refreshTick time.Duration, alertScanner AlertScanner) (*ProtectManager, error) {
	if refreshTick <= 0 {
		refreshTick = 10 * time.Second
	}

	if ttl <= 0 {
		ttl = 5 * time.Minute
	}

	if alertScanner == nil {
		return nil, fmt.Errorf("you must fill alert scanner")
	}

	ctx, cancel := context.WithCancel(context.Background())

	protectManager := &ProtectManager{
		ttl:           ttl,
		blocks:        make(map[string]time.Time),
		addChannel:    make(chan string, 1024),
		deleteChannel: make(chan string, 1024),
		refreshTick:   refreshTick,
		alertScanner:  alertScanner,
		ctx:           ctx,
		cancel:        cancel,
	}

	return protectManager, nil
}

func (m *ProtectManager) handleScan(ip string) {
	m.blocksMu.Lock()
	defer m.blocksMu.Unlock()

	now := time.Now()
	expireTime := now.Add(m.ttl)

	if _, exists := m.blocks[ip]; exists {
		return
	}

	if err := config.BlockIP(ip); err != nil {
		log.Printf("[ProtectManager] Block IP (%s) error: %v", ip, err)
		return
	}

	log.Printf("[ProtectManager] Blocked %s until %s", ip, expireTime.Format(time.RFC3339))

	m.blocks[ip] = expireTime
}

func (m *ProtectManager) handleExpiredBlock() {
	m.blocksMu.Lock()
	defer m.blocksMu.Unlock()

	now := time.Now()
	for ip, expireTime := range m.blocks {
		if now.After(expireTime) {
			if err := config.UnblockIP(ip); err != nil {
				log.Printf("[ProtectManager] Unblock IP (%s) error: %v", ip, err)
				continue
			}

			delete(m.blocks, ip)
			log.Printf("[ProtectManager] Unblocked %s (expired)", ip)
		}
	}
}

func (m *ProtectManager) handleDelete(ip string) {
	m.blocksMu.Lock()
	defer m.blocksMu.Unlock()

	expireTime, exists := m.blocks[ip]

	if !exists {
		return
	}

	if err := config.UnblockIP(ip); err != nil {
		log.Printf("[ProtectManager] Unblock IP (%s) error: %v", ip, err)
		return
	}

	delete(m.blocks, ip)
	log.Printf("[ProtectManager] Unblocked %s (manual, original-ttl %s)", ip, expireTime.Format(time.RFC3339))
}

func (m *ProtectManager) cleanupAll() {
	m.blocksMu.Lock()
	defer m.blocksMu.Unlock()

	for ip := range m.blocks {
		if err := config.UnblockIP(ip); err != nil {
			log.Printf("[ProtectManager] Unblock IP (%s) error: %v", ip, err)
			continue
		}

		log.Printf("[ProtectManager] Unblocked %s (cleanup)", ip)
	}

	m.blocks = make(map[string]time.Time)
}

func (m *ProtectManager) Add(ip string) {
	if ip == "" {

		return
	}

	select {
	case m.addChannel <- ip:
	default:
		log.Printf("[ProtectManager] Drop IP %s : Channel Full", ip)
	}
}

func (m *ProtectManager) Delete(ip string) {
	if ip == "" {
		return
	}

	select {
	case m.deleteChannel <- ip:
	default:
		log.Printf("[ProtectManager] Drop IP %s : Channel Full", ip)
	}
}

func (m *ProtectManager) SnapshotBlocks() []BlockedIP {
	m.blocksMu.RLock()
	defer m.blocksMu.RUnlock()

	snapshot := make([]BlockedIP, 0, len(m.blocks))
	for ip, expireAt := range m.blocks {
		snapshot = append(snapshot, BlockedIP{
			IP:          ip,
			ExpireAt:    expireAt,
			OriginalTTL: m.ttl,
		})
	}

	return snapshot
}

func (m *ProtectManager) Run(ctx context.Context) {
	readChannel := (<-chan string)(nil)
	readChannel = m.alertScanner.GetChannel()

	ticker := time.NewTicker(m.refreshTick)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			m.cleanupAll()
			return
		case ip := <-m.addChannel:
			m.handleScan(ip)
		case ip := <-m.deleteChannel:
			m.handleDelete(ip)
		case ip := <-readChannel:
			m.handleScan(ip)
		case <-ticker.C:
			m.handleExpiredBlock()
		}
	}
}

func (m *ProtectManager) Start() {
	m.alertScanner.StartScan()

	if m.ctx != nil {
		log.Println("[ProtectManager] Start scanning")
		go m.Run(m.ctx)
	}
}

func (m *ProtectManager) Stop() {
	m.alertScanner.StopScan()

	if m.cancel != nil {
		log.Println("[ProtectManager] Stop scanning")
		m.cancel()
	}
}
