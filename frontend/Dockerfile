# 완전한 빌드 이미지
FROM node:18

WORKDIR /workspace

COPY package*.json ./
RUN npm install

COPY . .

CMD ["npm", "run", "build"]