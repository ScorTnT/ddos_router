package protect_manager

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"
)

type Alert struct {
	EventType string `json:"event_type"`
	SrcIP     string `json:"src_ip"`
	DestIP    string `json:"dest_ip"`
}

// IP 를 선택해 반환하면 protect_manager 에서 차단.
// ok=false 이면 무시.
type FilterFn func(a *Alert) (ip string, ok bool)

// Reader 는 Snort 로그 파일을 tail 하며 IP 를 추출한다.
type Reader struct {
	path   string
	out    chan string
	filter FilterFn
}

// NewReader(path, buf, filter)docker
func NewReader(path string, buf int, filter FilterFn) *Reader {
	if buf <= 0 {
		buf = 128
	}
	return &Reader{
		path:   path,
		out:    make(chan string, buf),
		filter: filter,
	}
}

func (r *Reader) Output() <-chan string { return r.out }

// Run : ctx 가 취소될 때까지 tail + 파싱 루프
func (r *Reader) Run(ctx context.Context) {
	defer close(r.out)

	// 파일이 아직 없을 수도 있으므로 재시도 루프
	f, err := os.Open(r.path)
	for err != nil {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
			f, err = os.Open(r.path)
		}
	}
	defer f.Close()

	// 최신 로그 이후만 보려면 파일 끝으로 이동
	_, _ = f.Seek(0, io.SeekEnd)

	rd := bufio.NewReader(f)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			line, err := rd.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					time.Sleep(200 * time.Millisecond) // 새 로그 대기
					continue
				}
				log.Printf("[snortreader] read err: %v", err)
				time.Sleep(time.Second)
				continue
			}

			var a Alert
			if err := json.Unmarshal([]byte(line), &a); err != nil {
				continue // JSON 아니면 skip
			}

			if ip, ok := r.filter(&a); ok && ip != "" {
				select {
				case r.out <- ip: // 정상 전송
				default: // 채널이 가득 차면 드롭
					log.Printf("[snortreader] drop ip %s : chan full", ip)
				}
			}
		}
	}
}

/* ---------- 기본 Filter 예시 ---------- */

// BlockSrcIP : event_type==alert 인 라인의 src_ip
func BlockSrcIP(a *Alert) (string, bool) {
	if a.EventType == "alert" && a.SrcIP != "" {
		return a.SrcIP, true
	}
	return "", false
}

// BlockDestIP : dest_ip 사용
func BlockDestIP(a *Alert) (string, bool) {
	if a.EventType == "alert" && a.DestIP != "" {
		return a.DestIP, true
	}
	return "", false
}
