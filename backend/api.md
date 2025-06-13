# 내부 API 명세

| 구분        | 메소드 | 엔드포인트                 | 설명                                 |
| ----------- | ------ | -------------------------- | ------------------------------------ |
| Information | GET    | /information/neighbors     | 라우터 내 연결된 장비 리스트 리턴        |
|             | GET    | /information/connections   | 라우터 내, 외부로 이동하는 패킷 리스트 리턴 |
|             | GET    | /information               | 라우터 상세 정보 리턴                  |
| Config      | GET    | /config/wan                | 라우터 외부 인터페이스 설정 리턴         |
|             | POST   | /config/wan                | 라우터 외부 인터페이스 설정 저장         |
|             | GET    | /config/lan                | 라우터 내부 인터페이스 설정 리턴         |
|             | POST   | /config/lan                | 라우터 내부 인터페이스 설정 저장         |
| Protection  | GET    | /protection                | 프로텍션 로그(방어한 장비) 리스트 리턴   |
|             | POST   | /protection/block          | 아이피 수동 차단                       |
|             | POST   | /protection/unblock        | 아이피 수동 차단 해제                  |
| Auth        | POST   | /auth/login                | 계정 로그인 및 세션 코드 리턴           |
|             | POST   | /auth/logout               | 계정 로그아웃                          |