# OpenWrt 베이스 이미지
FROM openwrt/rootfs:x86-64-23.05.5

# /var/lock 디렉토리 생성 후 OpenWrt 패키지 업데이트 및 Git 설치
RUN mkdir -p /var/lock && opkg update && opkg install git git-http node curl nano

# npm 다운로드 및 설치
RUN curl -qL https://www.npmjs.com/install.sh | sh

# 작업 디렉토리 설정
WORKDIR /root

# Git 리포지토리 클론 명령 실행
RUN git clone https://github.com/ScorTnT/ddos_router

# frontend 디렉토리로 이동 및 npm 설치
WORKDIR /root/ddos_router/frontend
RUN npm install

# 포트 5173 개방
EXPOSE 5173

# 개발 서버 실행
CMD ["npm", "run", "dev"]