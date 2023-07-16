# atr

AT Protocol IPLD-CAR Repository toolkit

## Install

```
go install github.com/atscan/atr
```

## Examples

Scans the current directory (and its subdirectories) and prints information about all found repositories:
```bash
atr inspect
```

Return all objects from the repository:
```bash
atr show my-repo.car
```

Repositories can also be read via pipe:
```bash
curl 'https://bsky.social/xrpc/com.atproto.sync.getRepo?did=did:plc:524tuhdhh3m7li5gycdn6boe' | atr show
```