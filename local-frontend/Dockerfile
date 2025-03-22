# Node.js 기반 이미지 사용
FROM node:18

# 작업 디렉토리 설정
WORKDIR /app

# 의존성 파일이 없는 경우를 대비한 더미 생성
RUN touch package.json && touch package-lock.json

# 의존성 파일 복사 (없을 경우 더미 파일을 기반으로 설치)
COPY package*.json ./

# 의존성 설치
RUN npm install || echo "No dependencies to install"

# 소스 코드 복사
COPY . .

# 포트 노출
EXPOSE 5173

# 명령어 실행
CMD ["npm", "run", "dev"]