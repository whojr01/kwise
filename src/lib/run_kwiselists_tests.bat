set OLDGO=%GOPATH%
set GOPATH=e:\kwise_update
@rem go test -v kwiselists_test.go kwiselists.go
go test -run TestAddResult kwiselists_test.go kwisemap.go kwiselists.go kwiseFormats.go interface.go engine.go
set GOPATH=%OLDGO%
