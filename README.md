# How to build ?

## backend
```bash
> screen 1
cd ~/backend
./build.bat
./builds/backend
```
## frontend
choose option for how you want to excution . 

### run - realtime
```bash
> screen 2
cd ~/frontend 
npm install
npm run dev
```
### build and apply

PWD : `~/frontend`
```docker
docker build -t routerfront .

docker run --rm -v "${PWD}/dist:/workspace/dist" routerfront
```

# listen port
backend : 2024<br>

if you choose option run-realtime front you can connect this port to frontend .
frontend : 5173<br>