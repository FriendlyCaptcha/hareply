# hareply
`hareply` is a tool that replies to [HAProxy's `agent-check`](https://www.haproxy.com/documentation/haproxy-configuration-tutorials/service-reliability/health-checks/#agent-checks) that replies with the contents of a given file.

## CLI usage

```shell
# Print help
hareply -h

# Print help for serve command
hareply serve -h

# Print version information
hareply version

# Serve from port 8020 and a specific path.
hareply serve -f /some/path/to/agentstate -p 8020

# Serve from port 8442 and filepath `agentstate` (the defaults).
hareply serve
```

## As a library
`hareply` can also be used as a library, see the godoc for the `hareply` package.

## Error handling

* The response "`agentstate`" file is read on startup, if that fails the program will exit.
* The file is read again on any TCP connection, if that fails the last known file contents are used.
* If the value in the file is not a valid response for `agent-check`, the last valid response is returned instead.

## Release
Releases are built using `goreleaser`, see the [`goreleaser.yml`](./goreleaser.yml) file.

To mint a (test) release locally, install [goreleaser](https://goreleaser.com/install/) and run
```shell
goreleaser --snapshot --skip=publish --clean
```

## License
[MIT](./LICENSE.md) [🎶](https://suno.com/song/da6d4a83-1001-4694-8c28-648a6e8bad0a).