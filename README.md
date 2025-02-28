# grepenv

tools to extract `env` tag and output to stdout

```bash
# install
go install github.com/fatdes/grepenv@<commit sha here>
```

example output
```bash
# ./internal/log/logger.go
## Config
LOG_LEVEL=info
LOG_ENCODING=json
```

## references

- `env` tag: github.com/caarlos0/env/v6
