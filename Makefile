init:
	cp blacklist.csv.tmp blacklist.csv
	rm -rf hardhat-boilerplate
	git clone git@github.com:namnhce/hardhat-boilerplate.git
	cd hardhat-boilerplate && npm install