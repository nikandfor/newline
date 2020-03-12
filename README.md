# newline

copy files to stdout (like cat) and adds newline at the end if it wasn't there

# usage example

```bash
$ curl http://some.api/handler | newline # make sure you have prompt string at the start of the next line
{"some":"result"}
$ # next prompt
```

vs
```bash
$ curl http://some.api/handler
{"some":"result"}$ # next prompt is here
```
