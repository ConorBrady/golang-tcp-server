command -v go >/dev/null 2>&1 || {
    command -v brew >/dev/null 2>&1 || {
        ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
    }
    brew install go
}
go build server.go
