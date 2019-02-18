# promql-labelinjector

`promql-labelinjector` is a tool to inject/overwrite labels in a [PromQL](https://prometheus.io/docs/prometheus/latest/querying/basics/)-expression.

Closely related to the [PromAuthProxy](https://github.com/cherti/promauthproxy).

### manually

    # actually build and run
    git clone https://github.com/cherti/promql-labelinjector.git
    cd promql-labelinjector
    go get ./...
    go build promql-labelinjector.go
    ./promql-labelinjector -l job -v developer -e up


### automatically using go-toolchain

    go get -u "github.com/cherti/promql-labelinjector"
    ./promql-labelinjector -l job -v developer -e up

## Usage

`promql-labelinjector -l <labelname> -v <labelvalue> -e <promql-expression>` prints the promql-expression given with every metric having the specified label overwritten with the specified value or injected if it wasn't provided.

## License

This works is released under the [GNU General Public License v3](https://www.gnu.org/licenses/gpl-3.0.txt). You can find a copy of this license in the file `LICENSE` or at https://www.gnu.org/licenses/gpl-3.0.txt.

