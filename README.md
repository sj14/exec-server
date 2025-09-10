# exec-server

Execute commands sent via http.

## Example

### Run the server:

```
go run main.go
```

### Execute commands

By default, the repsonse arrives when the command finished:

```
➜ curl -X POST localhost:8080/exec -d "ls"    
LICENSE
main.go
README.md
```

For fire-and-forget add the respective header:

```
➜ curl -i -X POST localhost:8080/exec -d "sleep 10" -H "fire-and-forget: true"
```
