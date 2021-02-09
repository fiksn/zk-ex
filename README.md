# ZK-EX

[![tippin.me](https://badgen.net/badge/%E2%9A%A1%EF%B8%8Ftippin.me/@fiksn/F0918E)](https://tippin.me/@fiksn)

### Intro

Exclusive lock via Zookeeper.

Environment variables
```
ZKSERVERS="127.0.0.1:2181,127.0.0.2:2181,127.0.0.3:2181"
ZKLOCKPATH="/lock"
```

### Usage

Invoke like:
```
./zk-ex --lock "/bin/sh -c 'echo Lock acquired'" --nolock  "/bin/sh -c 'echo Lock not acquired'"
```
this way there will be no blocking (either lock or nolock command will get executed)

By omitting -nolock this is a blocking operation (it will wait until exclusive lock is acquired from zookeeper)
```
./zk-ex --lock "/bin/sh -c 'echo Finally acquired lock'" 
```

