# Sweep

[![build status](https://github.com/matteopolak/sweep/actions/workflows/build.yml/badge.svg)](.github/workflows/build.yml)
[![license](https://img.shields.io/github/license/matteopolak/sweep.svg)](LICENSE)

Sweep up heavy project files that are just collecting dust.

## Using a custom `sweep.toml`

Sweep will look for a `sweep.toml` file in the current directory. If one is not found, it will use the [default configuration](cmd/sweep/sweep.toml).

The format of a `sweep.toml` file is the following:

```toml
[[sweep]]
# folder to delete
folder = "node_modules"
# one of the files required to be in the same directory
predicate = ["package.json", "package-lock.json"]

[[sweep]]
folder = "target"
predicate = ["Cargo.toml"]
```
