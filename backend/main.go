package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/ScorTnT/ddos_router/backend/config"
	"github.com/ScorTnT/ddos_router/backend/controller"
	"github.com/ScorTnT/ddos_router/backend/database"

	"context"
	"time"

	pm "github.com/ScorTnT/ddos_router/backend/protect_manager"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	controller.HookHandler(app)

	var err error

	routerConfig, err := config.ReadConfig()

	if err != nil {
		log.Fatal(err)
		return
	}

	err = database.Init("secure_router.db", routerConfig)

	if err != nil {
		panic(err)
	}

	if err = config.InitFirewall(); err != nil {
		panic(err)
	}

	defer func() {
		if err := config.CleanupFirewall(); err != nil {
			fmt.Println(err)
		}
	}()

	/* Fiber, DB, Firewall 초기화 부분은 그대로 */

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Snort 로그 경로와 필터 정의
	const snortLog = "/var/log/snort/alert_json.txt"

	// protect_manager + Snort Reader 한번에 생성
	mgr := pm.NewWithReader(
		30*time.Minute, // TTL
		10*time.Second, // 만료 검사 주기
		snortLog,       // Snort JSON 로그 파일
		pm.BlockSrcIP,  // 필터 함수(예: src_ip 차단)
		256,            // Reader 채널 버퍼
	)
	go mgr.Run(ctx)

	/* 필요하다면 REST 등으로 수동 Protect 도 가능
	   e.g. curl POST /alert {"ip":"198.51.100.1"} */
	// router 설정 ...

	log.Fatal(app.Listen(":2024"))
}
