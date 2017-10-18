# Teatime

The workflow that we want users to have with our application is:

1. Start `teatime-sync` daemon.
2. Mark a directory as a synced "repository" destination
3. Either provide (connect to existing 'repository') or generate (create new
   repository) a key / hash / metadata that identifies the repository.
4. Daemon should poll marked directories, and connect to peers.
5. Daemons on peers communicate changes, and calculate diffs.
6. Changes should be communicated and applied on both ends.

### Peer to Peer Connection System

Peer discovery:
https://stackoverflow.com/questions/310607/peer-to-peer-methods-of-finding-peers

The provided key should contain a descriptive, unique ID of the
creator of the "repository". If a list of known peers (described later) is not
found, (i.e. IP addresses change), a centralized resource that should always be
up should be consulted to "fix" or "update" the list of known peers, and to mark
itself as a now-existing peer.

Connecting:
After discovering a peer, the two should connect (method of connection should be
decided based on implementation). Once connected, the peers should maintain a
keep-alive connection, and be able to communicate with each other without
needing to reconnect (unless connection is dropped).

Data Transfer:
The peers should be able to send arbitrary data amongst each other, but the data
should be verified on each end to be valid "Teatime" data.

### File Tracking System

There should be a way for the user to either whitelist or blacklist files. Ideas
are:

    - .sync-ignore file: Marks certain files or file extensions as untracked
    - User-added tracking: Only files that the user marks are tracked.

Tracked files are polled for differences.

### Diffing system

Files will be 'diffed' and a swapfile will be generated from the calculated
differences. The swapfile will be sent over the connection and propogated to all
peers.

### Merge Conflicts

In the event of 'merge conflict', there are several strategies to resolve the
issue:

    - User-interactive repair
    - Priority-heirarchy: Peers are in an ordered heiarchy, and 'higher-tiered'
      peers are deemed 'more correct'.
    - By relative timestamp (attempt to sync? This is heavily non-trivial).

