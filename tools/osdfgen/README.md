`osdfgen` can be used to build a CSV suitable for uploading to Pivotal's [OSDF Generator](http://osdf-generator.cfapps.io/static/index.html).
It determines licenses by sniffing the dependencies listed in `Gopkg.lock`.

Example:

```bash
go run osdfgen.go -p ../../ -o test.csv
```

The `-p` flag points at the project root and the `-o` flag is the place to put the output (stdout by default).
