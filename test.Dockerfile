FROM golang:1.15 AS test_runner

WORKDIR /test

COPY go.mod .
COPY go.sum .
RUN go mod download

RUN git config --global user.name "ign-cloudsim" && git config --global user.email "ign-cloudsim@test.org"

COPY . .

RUN go test -covermode=atomic -coverprofile=coverage.tx -v ./...

FROM sonarsource/sonar-scanner-cli AS sonar_scanner


