# 공식 Kali Linux 이미지를 기본 이미지로 사용
FROM kalilinux/kali-rolling

# 환경 변수 설정 (필요에 따라 추가 가능)
ENV DEBIAN_FRONTEND=noninteractive

# 기본 업데이트 및 필요한 패키지 설치 (예시로 기본 도구들 설치)
RUN apt-get update && \
    apt-get install -y \
    vim \
    curl && \
    rm -rf /var/lib/apt/lists/*

# 작업 디렉토리 설정 (필요에 따라)
WORKDIR /root

# 컨테이너 시작 시 실행할 기본 커맨드
CMD ["/bin/bash"]
