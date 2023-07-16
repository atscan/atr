# atr

AT Protocol IPLD-CAR Repository toolkit

## Install
```
go install github.com/atscan/atr
```

## Examples

```bash
# Scans the current directory (and its subdirectories) 
# and prints information  about all found repositories:
atr inspect

# Get all objects from the repository:
atr show my-repo.car

# Repositories can also be read via pipe:
curl -sL "https://xrpc.link/r/atproto.com" | atr show

# FYI xrpc.link is shortcut domain which redirecting to
# relevant /xrpc/... endpoints
```

## License

MIT