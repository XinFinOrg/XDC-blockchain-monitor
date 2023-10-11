all: build
	docker run -d -p 8000:8000 XinfinOrg/xinfin-monitor:latest

build:
	docker build -t liamlai/ai-server:latest .

run: build
	docker run -d -p 8000:8000 liamlai/ai-server:latest

build:
	docker build -t liamlai/ai-server:latest .

docker:
	docker build .

deploy:
	rsync -av ./ devnet3:~/xinfin-monitor/ && ssh devnet3 "make"