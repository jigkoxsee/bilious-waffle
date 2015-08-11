# bilious-waffle
Go web for fb-chat proxy

docker run -d -p 8080:8080 --name go-chat-proxy -e GO_PORT=8080 -e GO_HOST="http://chat.com" -e GO_TOKEN=TOKEN jigkoxsee/go-chat-proxy