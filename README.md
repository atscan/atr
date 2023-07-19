# atr

AT Protocol IPLD-CAR Repository toolkit

## Install
```
go install github.com/atscan/atr@latest
```

## Examples

```bash
# Scans the current directory (and its subdirectories) 
# and prints information  about all found repositories:
atr inspect
atr i  # you can use commands shortcuts

# Get all objects from the repository:
atr show my-repo.car

# filter by object type 
atr show -t post

# use jq query language
atr show -q .body.displayName

# use jmespath query language
atr show -q body.displayName

# Search with grep:
atr show -t post | grep love

# Repositories can also be read via pipe:
curl -sL "xrpc.link/r/atproto.com" | atr show

# FYI xrpc.link is shortcut domain which redirecting to
# relevant /xrpc/... endpoints
```

## License

MIT