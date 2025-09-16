package protection

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/ScorTnT/ddos_router/backend/internal/firewall"
)

type AlertScanner interface {
	StartScan()
	StopScan()
	GetChannel() <-chan string
}

type Manager struct {
	ttl                    time.Duration
	blocksMu               sync.RWMutex
	blocks                 map[string]BlockInfo // IP -> BlockInfo
	whitelistMu            sync.RWMutex
	whitelist              map[string]struct{}
	addChannel             chan string
	addPermanentChannel    chan string
	deleteChannel          chan string
	addWhitelistChannel    chan string
	deleteWhitelistChannel chan string
	refreshTick            time.Duration
	alertScanner           AlertScanner
	ctx                    context.Context
	cancel                 context.CancelFunc
}

type BlockInfo struct {
	ExpireTime  time.Time
	IsPermanent bool
}

type BlockedIP struct {
	IP          string        `json:"ip"`
	ExpireAt    time.Time     `json:"expire_at"`
	IsPermanent bool          `json:"is_permanent"`
	OriginalTTL time.Duration `json:"original_ttl"`
}

func NewProtectManager(ttl time.Duration, refreshTick time.Duration, alertScanner AlertScanner) (*Manager, error) {
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

	protectManager := &Manager{
		ttl:                    ttl,
		blocks:                 make(map[string]BlockInfo),
		addChannel:             make(chan string, 1024),
		addPermanentChannel:    make(chan string, 1024),
		deleteChannel:          make(chan string, 1024),
		whitelist:              make(map[string]struct{}),
		addWhitelistChannel:    make(chan string, 1024),
		deleteWhitelistChannel: make(chan string, 1024),
		refreshTick:            refreshTick,
		alertScanner:           alertScanner,
		ctx:                    ctx,
		cancel:                 cancel,
	}

	return protectManager, nil
}

func normalizeIP(ip string) (string, bool) {
	if ip == "" {
		return "", false
	}

	ip = strings.TrimSpace(ip)
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return "", false
	}
	// Canonical string form (handles IPv6 normalization as well)
	return parsed.String(), true
}

func (m *Manager) handleScan(ip string, isPermanent bool) {
	// Normalize incoming IP (covers Snort events and API enqueued adds)
	if nIP, ok := normalizeIP(ip); ok {
		ip = nIP
	} else {
		log.Printf("[ProtectManager] Ignore invalid IP on add: %q", ip)
		return
	}

	m.blocksMu.Lock()
	m.whitelistMu.Lock()
	defer m.blocksMu.Unlock()
	defer m.whitelistMu.Unlock()

	now := time.Now()
	expireTime := time.Time{}

	if !isPermanent {
		expireTime = now.Add(m.ttl)
	}

	if _, exists := m.blocks[ip]; exists {
		return
	}

	if _, exists := m.whitelist[ip]; exists {
		return
	}

	if err := firewall.BlockIP(ip); err != nil {
		// Idempotent: if it's already present in nft set, accept and proceed
		if list, gerr := firewall.GetBlockedIPs(); gerr == nil {
			found := false
			for _, e := range list {
				if e == ip {
					found = true
					break
				}
			}
			if !found {
				log.Printf("[ProtectManager] Block IP (%s) error: %v", ip, err)
				return
			}
			log.Printf("[ProtectManager] IP %s already present in firewall; syncing state", ip)
		} else {
			log.Printf("[ProtectManager] Block IP (%s) error and failed to verify state: %v", ip, err)
			return
		}
	} else {
		if isPermanent {
			log.Printf("[ProtectManager] Blocked %s permanently", ip)
		} else {
			log.Printf("[ProtectManager] Blocked %s until %s", ip, expireTime.Format(time.RFC3339))
		}
	}

	m.blocks[ip] = BlockInfo{
		ExpireTime:  expireTime,
		IsPermanent: isPermanent,
	}
}

func (m *Manager) handleExpiredBlock() {
	m.blocksMu.Lock()
	defer m.blocksMu.Unlock()

	now := time.Now()
	for ip, blockInfo := range m.blocks {
		if blockInfo.IsPermanent {
			continue
		}

		if now.After(blockInfo.ExpireTime) {
			if err := firewall.UnblockIP(ip); err != nil {
				log.Printf("[ProtectManager] Unblock IP (%s) error: %v", ip, err)
				continue
			}

			delete(m.blocks, ip)
			log.Printf("[ProtectManager] Unblocked %s (expired)", ip)
		}
	}
}

func (m *Manager) handleDelete(ip string) {
	// Normalize to match stored keys
	if nIP, ok := normalizeIP(ip); ok {
		ip = nIP
	} else {
		log.Printf("[ProtectManager] Ignore invalid IP on delete: %q", ip)
		return
	}

	m.blocksMu.Lock()
	defer m.blocksMu.Unlock()

	blockInfo, exists := m.blocks[ip]

	if !exists {
		return
	}

	if err := firewall.UnblockIP(ip); err != nil {
		log.Printf("[ProtectManager] Unblock IP (%s) error: %v", ip, err)
		return
	}

	delete(m.blocks, ip)
	log.Printf("[ProtectManager] Unblocked %s (manual, original-ttl %s)", ip, blockInfo.ExpireTime.Format(time.RFC3339))
}

func (m *Manager) handleAddWhitelist(ip string) {
	if nIP, ok := normalizeIP(ip); ok {
		ip = nIP
	} else {
		log.Printf("[ProtectManager] Ignore invalid IP on whitelist add: %q", ip)
		return
	}

	m.blocksMu.Lock()
	m.whitelistMu.Lock()
	defer m.whitelistMu.Unlock()
	defer m.blocksMu.Unlock()

	if _, exists := m.whitelist[ip]; !exists {
		m.whitelist[ip] = struct{}{}
		log.Printf("[ProtectManager] Add %s to whitelist", ip)
	}

	if blockInfo, exists := m.blocks[ip]; exists {
		if err := firewall.UnblockIP(ip); err != nil {
			log.Printf("[ProtectManager] Unblock IP (%s) on whitelist add error: %v", ip, err)
			return
		}
		delete(m.blocks, ip)
		log.Printf("[ProtectManager] Unblocked %s (added to whitelist, original-ttl %s)", ip, blockInfo.ExpireTime.Format(time.RFC3339))
	}

	m.whitelist[ip] = struct{}{}
	log.Printf("[ProtectManager] Add %s to whitelist", ip)
}

func (m *Manager) handleDeleteWhitelist(ip string) {
	if nIP, ok := normalizeIP(ip); ok {
		ip = nIP
	} else {
		log.Printf("[ProtectManager] Ignore invalid IP on whitelist delete: %q", ip)
		return
	}

	m.whitelistMu.Lock()
	defer m.whitelistMu.Unlock()

	if _, exists := m.whitelist[ip]; !exists {
		return
	}

	delete(m.whitelist, ip)
	log.Printf("[ProtectManager] Delete %s from whitelist", ip)
}

func (m *Manager) cleanupAll() {
	m.blocksMu.Lock()
	defer m.blocksMu.Unlock()

	for ip := range m.blocks {
		if err := firewall.UnblockIP(ip); err != nil {
			log.Printf("[ProtectManager] Unblock IP (%s) error: %v", ip, err)
			continue
		}

		log.Printf("[ProtectManager] Unblocked %s (cleanup)", ip)
	}

	m.blocks = make(map[string]BlockInfo)
}

func (m *Manager) Add(ip string, isPermanent bool) {
	if nIP, ok := normalizeIP(ip); ok {
		ip = nIP
	} else {
		log.Printf("[ProtectManager] Ignore invalid IP on add enqueue: %q", ip)
		return
	}

	if isPermanent {
		select {
		case m.addPermanentChannel <- ip:
		default:
			log.Printf("[ProtectManager] Drop Permanent IP %s : Channel Full", ip)
		}
	} else {
		select {
		case m.addChannel <- ip:
		default:
			log.Printf("[ProtectManager] Drop IP %s : Channel Full", ip)
		}
	}
}

func (m *Manager) Delete(ip string) {
	if nIP, ok := normalizeIP(ip); ok {
		ip = nIP
	} else {
		log.Printf("[ProtectManager] Ignore invalid IP on delete enqueue: %q", ip)
		return
	}

	select {
	case m.deleteChannel <- ip:
	default:
		log.Printf("[ProtectManager] UnDrop IP %s : Channel Full", ip)
	}
}

func (m *Manager) AddWhitelist(ip string) {
	if nIP, ok := normalizeIP(ip); ok {
		ip = nIP
	} else {
		log.Printf("[ProtectManager] Ignore invalid IP on whitelist add enqueue: %q", ip)
		return
	}

	select {
	case m.addWhitelistChannel <- ip:
	default:
		log.Printf("[ProtectManager] Add Ip to Whitelist %s : Channel Full", ip)
	}
}

func (m *Manager) DeleteWhitelist(ip string) {
	if nIP, ok := normalizeIP(ip); ok {
		ip = nIP
	} else {
		log.Printf("[ProtectManager] Ignore invalid IP on whitelist delete enqueue: %q", ip)
		return
	}

	select {
	case m.deleteWhitelistChannel <- ip:
	default:
		log.Printf("[ProtectManager] Delete Ip from Whitelist %s : Channel Full", ip)
	}
}

func (m *Manager) SnapshotBlocks() []BlockedIP {
	m.blocksMu.RLock()
	defer m.blocksMu.RUnlock()

	snapshot := make([]BlockedIP, 0, len(m.blocks))
	for ip, blockInfo := range m.blocks {
		snapshot = append(snapshot, BlockedIP{
			IP:          ip,
			ExpireAt:    blockInfo.ExpireTime,
			IsPermanent: blockInfo.IsPermanent,
			OriginalTTL: m.ttl,
		})
	}

	return snapshot
}

func (m *Manager) SnapshotWhitelist() []string {
	m.whitelistMu.RLock()
	defer m.whitelistMu.RUnlock()

	snapshot := make([]string, 0, len(m.whitelist))
	for ip := range m.whitelist {
		snapshot = append(snapshot, ip)
	}

	return snapshot
}

func (m *Manager) Run(ctx context.Context) {
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
			m.handleScan(ip, false)
		case ip := <-m.addPermanentChannel:
			m.handleScan(ip, true)
		case ip := <-m.deleteChannel:
			m.handleDelete(ip)
		case ip := <-readChannel:
			m.handleScan(ip, false)
		case ip := <-m.addWhitelistChannel:
			m.handleAddWhitelist(ip)
		case ip := <-m.deleteWhitelistChannel:
			m.handleDeleteWhitelist(ip)
		case <-ticker.C:
			m.handleExpiredBlock()
		}
	}
}

func (m *Manager) Start() {
	m.alertScanner.StartScan()

	if m.ctx != nil {
		// Sync current firewall state into memory before starting event loop
		m.syncFromFirewall()
		log.Println("[ProtectManager] Start scanning")
		go m.Run(m.ctx)
	}
}

func (m *Manager) Stop() {
	m.alertScanner.StopScan()

	if m.cancel != nil {
		log.Println("[ProtectManager] Stop scanning")
		m.cancel()
	}
}

func (m *Manager) syncFromFirewall() {
	list, err := firewall.GetBlockedIPs()
	if err != nil {
		log.Printf("[ProtectManager] Failed to sync from firewall: %v", err)
		return
	}

	if len(list) == 0 {
		return
	}

	now := time.Now()
	expire := now.Add(m.ttl)

	m.blocksMu.Lock()
	m.whitelistMu.Lock()
	defer m.whitelistMu.Unlock()
	defer m.blocksMu.Unlock()

	count := 0
	for _, ip := range list {
		if n, ok := normalizeIP(ip); ok {
			ip = n
		} else {
			continue
		}

		if _, w := m.whitelist[ip]; w {
			continue
		}
		if _, ex := m.blocks[ip]; ex {
			continue
		}
		m.blocks[ip] = BlockInfo{
			ExpireTime:  expire,
			IsPermanent: false,
		}
		count++
	}

	if count > 0 {
		log.Printf("[ProtectManager] Synced %d IP(s) from firewall", count)
	}
}
