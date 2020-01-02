IMAGE=lukapiske/aloha
RELEASE=`cat ./VERSION`
IMAGE_TAG=${IMAGE}:${RELEASE}

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  .



