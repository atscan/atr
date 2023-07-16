# atr

AT Protocol IPLD-CAR Repository toolkit

## Install
```
go install github.com/atscan/atr
```

## Examples

```bash
# Scans the current directory (and its subdirectories) and prints information about all found repositories:
atr inspect

# Get all objects from the repository:
atr show my-repo.car

# Repositories can also be read via pipe:
curl "https://bsky.social/xrpc/com.atproto.sync.getRepo?did=did:plc:524tuhdhh3m7li5gycdn6boe" | atr show
```

## License

MIT