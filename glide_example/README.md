
## how to reproduce

configure to run glide

    $ cd path/to/glide_sample
    $ export GOPATH=$(pwd)
    $ export PATH="$(pwd)/bin:$PATH"
    $ cd src/glide_sample/
    $ glide install

confirm glide process

    $ ps -ax | grpe glide

