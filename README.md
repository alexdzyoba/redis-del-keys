redis-del-keys
==============

Delete keys in Redis by a pattern. Deletion is performed by doing [`SCAN`][scan]
and then invoking DEL in a [pipeline]. This allows you to delete the keys in a
non-blocking manner.

## Usage

Simple invocation requires pattern argument:

    $ redis-del-keys 
    pattern argument is required

    redis-iter-del iterates over redis keys with SCAN matched by pattern and then
    DEL the keys in pipelined commands

    Usage of redis-del-keys:
      -a, --addrs strings    Redis addrs, comma separated for cluster (default [:6379])
      -b, --batch int        Batch size for pipelined commands (default 10)
      -c, --count int        Count for SCAN command (default 10)
      -d, --dryrun           Dry run
      -p, --pattern string   Pattern to delete

This prevents you from accidentally deleting all of the keys.

Real example:

    $ redis-del-keys -p 'a*'
    iterated over 10 keys

It prints how many keys it has iterated. Running this command again will iterate
over 0 keys because they were deleted:

    $ redis-del-keys -p 'a*'
    iterated over 0 keys

The command supports **dry run mode** via `-d` or `--dryrun` option. It is safe
and allows you to estimate how many keys will be deleted:

    $ redis-del-keys -d -p '*'
    iterated over 240 keys


[scan]: https://redis.io/commands/scan
[pipeline]: https://redis.io/topics/pipelining
