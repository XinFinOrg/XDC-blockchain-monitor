all: build
	docker run -d xinfinorg/xdc-blockchain-monitor XDC-blockchain-monitor

build:
	docker build -t xinfinorg/xdc-blockchain-monitor:latest .

deploy:
	rsync -av ./ devnet3:~/xdc-blockchain-monitor/ && ssh devnet3 "cd xdc-blockchain-monitor && make"