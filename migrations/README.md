To create a migration

```
cd migrations
goose create galleries sql
```


Because we are developing solo, we can remove the timestamp
```
goose fix 
```


```
-- +goose Up
-- +goose StatementBegin
CREATE TABLE galleries (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users (id),
  title TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE galleries;
-- +goose StatementEnd
```

```
cd ..
code models/gallery.go
```

```
package models

type Gallery struct {
	ID     int
	UserID int
	Title  string
}
```