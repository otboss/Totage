
Ensure go paths are set up correctly, On MacOS or Linux add following lines to `.bashrc` if not already added
```
export GOROOT=$HOME/go
export GOPATH=$HOME/go/pkg
export PATH=$GOROOT/bin:$PATH
export PATH=$GOROOT/pkg/bin:$PATH
```

Install protobuf related packages if not already available:
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/bufbuild/buf/cmd/buf@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

Install protobuf related packages to generate typescript if not already available:
```
npm install -g @bufbuild/protoc-gen-es @connectrpc/protoc-gen-connect-es
```

In the project root directory, generate Proto stubs using the following command:
```
buf generate --include-imports
```

In the project root directory, generate models using the following command:
```
sqlc generate -f ./sqlc.yml
```
<!-- ```
(cd ./repository/sqlc && sqlc generate)
``` -->

Setup Kafka using the `kafka.yml` file located in the project root directory:
```
```

Download k3s from the official repository for Kubernetes support:
```
curl -sL $(curl -s https://api.github.com/repos/k3s-io/k3s/releases/latest | grep "browser_download_url" | cut -d '"' -f 4) -o k3s
```

Run project using the run.sh file in the project root directory:
```
./run.sh
```