# protocenv

manage protoc versions.

## Installation

```
go get -u github.com/garupanojisan/protocenv
```
## PATH Configuration

protocenv creates symbolic link at $HOME/.protocenv/bin to global specified version's bin,
so please set PATH as below.
```
export PATH=$HOME/.protocenv/bin:$PATH
```

## Usage

```
# show list of all available versions
protocenv install -l

# install specified version
protocenv install <version>

# set global version
protocenv global <version>

# set local version
protocenv local <version>
```
