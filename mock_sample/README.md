# gomock sample with aws-sdk-go

Run the mocking test with aws-sdk-go

```bash
$ make test
mockgen -source path/to/pkg/mod/github.com/aws/aws-sdk-go@v1.29.17/service/s3/s3iface/interface.go -package main -destination ./mockS3.go
go test -v .
=== RUN   TestGetObject
2020/03/05 17:20:08 GetObjectOutput: {
  Body: buffer(%!p(ioutil.nopCloser={0xc0000b4f40}))
}
2020/03/05 17:20:08 wrote file: single-line.txt
2020/03/05 17:20:08 GetObjectOutput: {
  Body: buffer(%!p(ioutil.nopCloser={0xc0000b4f60}))
}
2020/03/05 17:20:08 wrote file: multiple-line.txt
--- PASS: TestGetObject (0.00s)
PASS
ok  	example.com/mock_sample	0.003s
```

Some text files are created by [main_test.go](./main_test.go).

```bash
$ cat single-line.txt 
dummy

$ cat multiple-line.txt 
one
two
three
```
