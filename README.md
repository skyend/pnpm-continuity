# Pnpm inter continuity

Mirroring project package.json
```shell
go run inter-continuity/mirroring.go
```
Will be created `npm-packages/` in your working directory


Publish mirrored modules to private repository
```shell
go run inter-continuity/publish.go
```
Do change `npm registry url` and `npm login` before execute `publish.go`
