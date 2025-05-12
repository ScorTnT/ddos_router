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

type FilterFunc func(a *Alert) (ip string, ok bool)

type SnortScanner struct {
	logPath    string
	outChannel chan string
	filter     FilterFunc
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewSnortScanner(logPath string, bufferSize int, filter FilterFunc) *SnortScanner {
	if bufferSize <= 0 {
		bufferSize = 128
	}

	if filter == nil {
		filter = DefaultFilter
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &SnortScanner{
		logPath:    logPath,
		outChannel: make(chan string, bufferSize),
		filter:     filter,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (s *SnortScanner) GetChannel() <-chan string { return s.outChannel }

func (s *SnortScanner) Run(ctx context.Context) {
	defer close(s.outChannel)

	f, err := os.Open(s.logPath)
	for err != nil {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
			f, err = os.Open(s.logPath)
		}
	}

	defer f.Close()

	_, _ = f.Seek(0, io.SeekEnd)

	logReader := bufio.NewReader(f)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			line, err := logReader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					time.Sleep(200 * time.Millisecond)
					continue
				}
				log.Printf("[SnortScanner] Read Error: %v", err)
				time.Sleep(time.Second)
				continue
			}

			var alert Alert
			if err := json.Unmarshal([]byte(line), &alert); err != nil {
				continue
			}

			if ip, ok := s.filter(&alert); ok && ip != "" {
				select {
				case s.outChannel <- ip:
				default:
					log.Printf("[SnortScanner] Drop IP %s : Channel Full", ip)
				}
			}
		}
	}
}

func (s *SnortScanner) StartScan() {
	if s.logPath == "" {
		return
	}

	go s.Run(s.ctx)
}

func (s *SnortScanner) StopScan() {
	if s.logPath == "" {
		return
	}

	if s.cancel != nil {
		s.cancel()
	}
}
