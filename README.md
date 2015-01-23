### gochunk


#### chop

Chop a large file into smaller chunks and gzip them. The chunks
are named after the sha1sum of its (gzipped) contents. The sha1s
are printed in order to stdout. You should probably pipe that
into a file.

```sh
    $ ./gochunk chop bigfile.bin | tee bigfile.bin.manifest
    5ca1ac2a781798dbeb8d1030278783788a1ac057    bigfile.bin
    aec7fd9e0af9ac429ab86c6aafef36461ea33fda    bigfile.bin
    d85ef40572702425801ac3ea1bb06212dce12716    bigfile.bin
    d1a4d8858b90faf86b4b9c46866e99acc65bd25e    bigfile.bin
    3b4a3ae42991722caf167c74007392aa808e3187    bigfile.bin
    9ed0bea1f0d9b3fb52a12ef1d2352ccc6b050e7a    bigfile.bin
    8e71c4b54532842cdcf48f8cee380e1d1bb882d2    bigfile.bin
```

You'll see a directory `chunks`. This contains the chopped up
bigfile.bin. It spreads the files across many subdirectories,
by sharding them based on the first few characters of the
filename.

```sh
    $ find chunks/ -type f
    chunks/5c/a1/ac/5ca1ac2a781798dbeb8d1030278783788a1ac057
    chunks/ae/c7/fd/aec7fd9e0af9ac429ab86c6aafef36461ea33fda
    chunks/d8/5e/f4/d85ef40572702425801ac3ea1bb06212dce12716
    chunks/d1/a4/d8/d1a4d8858b90faf86b4b9c46866e99acc65bd25e
    chunks/3b/4a/3a/3b4a3ae42991722caf167c74007392aa808e3187
    chunks/9e/d0/be/9ed0bea1f0d9b3fb52a12ef1d2352ccc6b050e7a
    chunks/8e/71/c4/8e71c4b54532842cdcf48f8cee380e1d1bb882d2
```

#### serve

Run `gochunk server` to start a static webserver serving out your
chunks to your peers.

```sh
    $ ./gochunk server
    Listening on http://localhost:9999
```


