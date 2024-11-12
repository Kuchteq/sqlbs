# sqlbs - come up with plausible bullshit for your sql database

The front-end guys complain that they don't really know how the UI's going to look like once more data comes in. You're creating a backend and you don't know if your SQL queries really return what they are supposed to? *Bsql* solves that issue by generating a bunch of insert statements based on sql schema comments specifying one of the many predefined collections.

### In practice
Annotate your schemas in the following way:

```
CREATE TABLE members (
        id INTEGER NOT NULL PRIMARY KEY,
        username TEXT NOT NULL, --bs: username
        displayName TEXT NOT NULL, --bs: fullname
);
```

now run **sqlbs** which prints insert statements that you can then pipe to the db tool of your choice.
