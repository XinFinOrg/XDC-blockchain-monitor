all: build
	docker stop xdc-blockchain-monitor || true
	docker rm xdc-blockchain-monitor || true
	docker run -p 8080:8080 --name xdc-blockchain-monitor --restart=always -d xinfinorg/xdc-blockchain-monitor XDC-blockchain-monitor

build:
	docker build -t xinfinorg/xdc-blockchain-monitor:latest .

deploy:
	rsync -av ./ devnet3:~/xdc-blockchain-monitor/ && ssh devnet3 "cd xdc-blockchain-monitor && make"