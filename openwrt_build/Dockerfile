# OpenWrt 베이스 이미지
FROM openwrt/rootfs:x86-64-23.05.5

# /var/lock 디렉토리 생성 후 OpenWrt 패키지 업데이트 및 Git 설치
RUN mkdir -p /var/lock && opkg update && opkg install git git-http node curl nano python3-pip bash dos2unix

# supervisor 설치
RUN pip3 install supervisor

# npm 다운로드 및 설치
RUN curl -qL https://www.npmjs.com/install.sh | sh

# 작업 디렉토리 설정
WORKDIR /root

# Git 리포지토리 클론 명령 실행
RUN git clone https://github.com/ScorTnT/ddos_router

# frontend 디렉토리로 이동 및 npm 설치
WORKDIR /root/ddos_router/frontend
RUN npm install

# backend 디렉토리 이동
WORKDIR /root/ddos_router/backend

# build 파일에 실행 권한 부여
RUN chmod +x build

# Windows 스타일 줄바꿈 형식 변환
RUN dos2unix build

# build 실행 파일 생성 (bash 사용)
RUN ls -l && bash build

# backend 실행파일에 실행 권한 부여
WORKDIR /root/ddos_router/backend/builds/
RUN chmod +x backend

# 로그 파일 저장할 디렉토리 생성
RUN mkdir -p /var/log/

# supervisord.conf 파일을 컨테이너로 복사
COPY supervisord.conf /etc/supervisord.conf

# 포트 개방 (프론트엔드: 5173, 백엔드: 2024)
EXPOSE 5173 2024

# supervisor로 두 프로세스 실행
CMD ["supervisord", "-c", "/etc/supervisord.conf"]
