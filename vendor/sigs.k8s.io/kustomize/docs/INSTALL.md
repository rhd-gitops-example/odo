[release page]: https://github.com/kubernetes-sigs/kustomize/releases
[Go]: https://golang.org

<<<<<<< HEAD
## Installation

On macOS, you can install kustomize with Homebrew package
manager:

    brew install kustomize

For all operating systems, download a binary from the
[release page].

Or try this to grab the latest official release
using the command line:
=======
# Installation

Binaries at various versions for linux, macOs and Windows
are available on the [release page].

Or...

## Quickly curl the latest
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)

```
opsys=linux  # or darwin, or windows
curl -s https://api.github.com/repos/kubernetes-sigs/kustomize/releases/latest |\
  grep browser_download |\
  grep $opsys |\
  cut -d '"' -f 4 |\
  xargs curl -O -L
mv kustomize_*_${opsys}_amd64 kustomize
chmod u+x kustomize
```

<<<<<<< HEAD
To install from head with [Go] v1.10.1 or higher:

<!-- @installkustomize @test -->
```
go get sigs.k8s.io/kustomize
```
=======
## Get and install source for a particular release

For example
```
# Omit the @v3.2.1 to get the default for major version 3
GO111MODULE=on go get sigs.k8s.io/kustomize/kustomize/v3@v3.2.1
```

Use of `GO111MODULE=on` shouldn't be necessary
with [Go v1.13](https://golang.org/doc/go1.13#modules).

### Other methods

#### macOS

```
brew install kustomize
```

#### windows

```
choco install kustomize
```

For support on the chocolatey package
and prior releases, see:
- [Choco Package](https://chocolatey.org/packages/kustomize)
- [Package Source](https://github.com/kenmaglio/choco-kustomize)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
