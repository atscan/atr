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
curl -sL "https://enoki.us-east.host.bsky.network/xrpc/com.atproto.sync.getRepo?did=did:plc:ewvi7nxzyoun6zhxrhs64oiz" | atr show
curl -sL "xrpc.link/r/atproto.com" | atr show

# FYI xrpc.link is shortcut domain which redirecting to
# relevant /xrpc/... endpoints
```

Example output:
```bash
~> curl -sL "https://enoki.us-east.host.bsky.network/xrpc/com.atproto.sync.getRepo?did=did:plc:ewvi7nxzyoun6zhxrhs64oiz" | atr i

(pipe):
  DID: did:plc:ewvi7nxzyoun6zhxrhs64oiz  Repo Version: 3
  Head: bafyreiapeyhetsiuz6jcybcmgxw7k45imuu22shmbhybccq4vef57fc4mi
  Sig: d03ed428c77a9e4964560f575426f81ed694861c93372b866d77b3b1f814a39b4d27033fa86593e194dfe07278a96aa1a9b76f1b2f7ff8d4313130b557177fb1
  Size: 566 kB  Blocks: 1,837  Commits: 1  Objects: 1,434
  Profile:
    Display Name: AT Protocol Developers
    Description: Social networking technology created by Bluesky.

      Developer-focused account. Follow @bsky.app for general announcements!

      Bluesky API docs: docs.bsky.app
      AT Protocol specs: atproto.com
  Collections:
    app.bsky.actor.profile: 1
    app.bsky.feed.generator: 2
    app.bsky.feed.like: 686
    app.bsky.feed.post: 150
    app.bsky.feed.repost: 541
    app.bsky.feed.threadgate: 1
    app.bsky.graph.follow: 21
    app.bsky.graph.list: 1
    app.bsky.graph.listitem: 29
    app.bsky.graph.starterpack: 1
    chat.bsky.actor.declaration: 1
  Last 5 commits:
    bafyreiapeyhetsiuz6jcybcmgxw7k45imuu22shmbhybccq4vef57fc4mi
```

## License

MIT
