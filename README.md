NOTE
----

Note creating program

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
