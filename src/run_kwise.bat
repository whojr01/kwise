@echo Don't for get to set GOPATH to: set GOPATH=E:\kwise
set GOPATH=E:\kwise
@REM go run main.go --pageDescription lfs_blfs_page.kwi --dataFile blfs_dbus.html --trace
@REM go run main.go --pageDescription simple_page.kwi --dataFile simple_page_data.html --trace
@REM go run main.go --pageDescription simple_page.kwi --dataFile blfs_dbus.html --trace

go run main.go --pageDescription simple_page.kwi --dataFile simple_page_data.html

@REM go run main.go --pageDescription simple_collectnext_page.kwi --dataFile simple_collectnext_data.html