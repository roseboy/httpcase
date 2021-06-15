
build:
	go build -o hc main.go

install:
	@go build -o hc main.go
	@mkdir -p /usr/local/httpcase
	@cp hc /usr/local/httpcase/
	@cp README.md /usr/local/httpcase/
	@cp LICENSE /usr/local/httpcase/
	@cp -rf plugin-js /usr/local/httpcase/
	@ln -snf /usr/local/httpcase/hc /usr/local/bin/hc
	hc version

rl:
	goreleaser --snapshot --skip-publish --rm-dist

twine:
	cd pip-install-httpcase \
	&& python setup.py sdist \
	&& python setup.py bdist_wheel --universal \
	&& twine upload dist/*