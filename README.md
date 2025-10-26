Tutorial link: https://timeweb.cloud/tutorials/docker/ispolzovanie-selenium-s-chrome-v-docker

Create .env file and fill it with data using example.env\
Run: go run .\cmd\main.go\
Build: docker compose up -d --build\
Open testing browser from docker - http://localhost:7900  \
Change VNC password steps:
1) docker exec -it selenium-container bash
2) x11vnc -storepasswd

