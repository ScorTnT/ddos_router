package protect_manager

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ScorTnT/ddos_router/backend/config"
)

// Manager 는 IP 차단/해제 + Snort 로그 모니터링까지 담당한다.
type Manager struct {
	ttl       time.Duration        // IP 차단 유지 기간
	mu        sync.Mutex           // blocks 맵 보호용 뮤텍스
	blocks    map[string]time.Time // ip -> 만료시각
	addCh     chan string          // 수동/Reader 양쪽에서 들어오는 IP
	tickerDur time.Duration        // 만료 검사용 주기

	// Snort Reader (선택)
	reader *Reader
}

// NewManager : Snort 연동이 필요 없는 경우
func NewManager(ttl, tick time.Duration) *Manager {
	if tick <= 0 {
		tick = 10 * time.Second
	}
	return &Manager{
		ttl:       ttl,
		blocks:    make(map[string]time.Time),
		addCh:     make(chan string, 1024),
		tickerDur: tick,
	}
}

// NewWithReader : Snort JSON 로그를 함께 모니터링하고 싶을 때 사용하는 생성자
func (m *Manager) handleNewIP(ip string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	exp := now.Add(m.ttl)

	// 이미 차단 중이면 만료시각만 갱신
	if _, exists := m.blocks[ip]; exists {
		m.blocks[ip] = exp
		return
	}

	// 신규 차단
	if err := config.BlockIP(ip); err != nil {
		log.Printf("[protect_manager] BlockIP(%s) error: %v", ip, err)
		return
	}
	log.Printf("[protect_manager] Blocked %s until %s", ip, exp.Format(time.RFC3339))
	m.blocks[ip] = exp
}

// expireBlocks : 주기적으로 만료 여부 확인
func (m *Manager) expireBlocks() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for ip, exp := range m.blocks {
		if now.After(exp) {
			if err := config.UnblockIP(ip); err != nil {
				log.Printf("[protect_manager] UnblockIP(%s) error: %v", ip, err)
				continue
			}
			delete(m.blocks, ip)
			log.Printf("[protect_manager] Unblocked %s (expired)", ip)
		}
	}
}

// cleanupAll : 애플리케이션 종료 시 모든 IP 해제
func (m *Manager) cleanupAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for ip := range m.blocks {
		if err := config.UnblockIP(ip); err != nil {
			log.Printf("[protect_manager] UnblockIP(%s) error: %v", ip, err)
			continue
		}
		log.Printf("[protect_manager] Unblocked %s (shutdown)", ip)
	}
	m.blocks = make(map[string]time.Time) // 깔끔히 초기화
}

// NewWithReader : Snort JSON 로그를 함께 모니터링하고 싶을 때 사용하는 생성자
func NewWithReader(ttl, tick time.Duration, logPath string, filter FilterFn, buf int) *Manager {
	m := NewManager(ttl, tick) // 기본 매니저 생성
	m.reader = NewReader(logPath, buf, filter)
	return m
}

// Add : 외부에서 수동으로 IP 차단을 요청할 때 사용
func (m *Manager) Add(ip string) {
	if ip == "" {
		return
	}
	select {
	case m.addCh <- ip:
	default: // 채널이 가득 차면 드롭
		log.Printf("[protect_manager] drop ip %s : chan full", ip)
	}
}

// Run : 컨텍스트가 종료될 때까지 메인 루프 실행
//  1. 수동 IP(addCh)       2) Snort Reader 결과
//  3. 주기적 만료 검사     4) 종료 시 모든 IP 해제
func (m *Manager) Run(ctx context.Context) {
	// Reader가 있으면 먼저 실행
	if m.reader != nil {
		go m.reader.Run(ctx)
	}

	readCh := (<-chan string)(nil)
	if m.reader != nil {
		readCh = m.reader.Output()
	}

	ticker := time.NewTicker(m.tickerDur)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			m.cleanupAll()
			return

		case ip := <-m.addCh: // 수동 추가
			m.handleNewIP(ip)

		case ip := <-readCh: // Snort Reader 결과
			m.handleNewIP(ip)

		case <-ticker.C: // 주기적 만료 검사
			m.expireBlocks()
		}
	}
}
