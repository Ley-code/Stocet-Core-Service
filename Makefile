.PHONY: build clean package

# Build a statically-linked binary for AWS Lambda (Amazon Linux 2023 / Graviton).
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -tags lambda.norpc \
		-o bootstrap cmd/server/main.go

# Package the binary into a deployment zip for Lambda upload.
package: build
	zip deployment.zip bootstrap

clean:
	rm -f bootstrap deployment.zip
