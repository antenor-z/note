NOTE
----

Note creating program

## Development
Use the same instruction for deploying, but with the compose
files ending with "-dev".

You can also run without Docker. On the note-back run:
```
go mod download
go build
./note
```

On the note-front, you can use a simple Python server
```
python3 -m http.server 3000
```

## Deploying

- Inside note-back folder, create the following files:

auth.toml:
```
username="changeme"
password="changeme"
```

config.toml:
```
domain="http://localhost:3003"
debugmode=true
```

- Create empty file "db.db".

- Build and up compose(-dev).yml on note-back and note-front.
