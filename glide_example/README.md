
## how to reproduce

configure to run glide

    $ cd path/to/glide_sample
    $ export GOPATH=$(pwd)
    $ export PATH="$(pwd)/bin:$PATH"
    $ go get github.com/Masterminds/glide
    $ cd src/glide_sample/
    $ glide install

confirm glide process

    $ ps -ax | grpe glide

