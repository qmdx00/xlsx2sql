# xlsx2sql

Parse xlsx file and convert content to sql statements, support batch insert.

## Installation

```bash
go get github.com/qmdx00/xlsx2sql
```

## Usage


Example test.xlsx file:

| id  | name   |
| --- | ------ |
| 1   | test01 |
| 2   | test02 |
| 3   | test03 |

### Single insert

```go
xlsx.New("test.xlsx", xlsx.SQLMode(xlsx.ModeSingle), xlsx.TableName("test")).
  Mapped("id", 0).          // get id from xlsx column 0
  Mapped("name", 1).        // get name from xlsx column 1
  Columns("key", "value").  // set key with value
  Build()
```

Output:
```SQL
INSERT INTO `test` (`id`, `name`, `key`) VALUES ('1', 'test01', 'value');
INSERT INTO `test` (`id`, `name`, `key`) VALUES ('2', 'test02', 'value');
INSERT INTO `test` (`id`, `name`, `key`) VALUES ('3', 'test03', 'value');
```

### Batch insert

```go
xlsx.New("test.xlsx", xlsx.SQLMode(xlsx.ModeBatch), xlsx.BatchSize(100), xlsx.TableName("test")).
  Mapped("id", 0).          // get id from xlsx column 0
  Mapped("name", 1).        // get name from xlsx column 1
  Columns("key", "value").  // set key with value
  Build()
```

Output:
```SQL
INSERT INTO `test` (`id`, `name`, `key`) VALUES
('1', 'test01', 'value'),
('2', 'test02', 'value'),
('3', 'test03', 'value');
```

### Custom value processor

```go
xlsx.New("test.xlsx", xlsx.SQLMode(xlsx.ModeBatch), xlsx.BatchSize(100), xlsx.TableName("test")).
  Mapped("id", 0).          // get id from xlsx column 0
  Mapped("name", 1).        // get name from xlsx column 1
  Columns("key", "value").  // set key with value
  Valuer("id", func(v string) string {
    id, _ := strconv.Atoi(v)
    return fmt.Sprintf("%d", id+1)
  }).Build()
```

Output:
```SQL
INSERT INTO `test` (`id`, `name`, `key`) VALUES
('2', 'test01', 'value'),
('3', 'test02', 'value'),
('4', 'test03', 'value');
```

### Mappping columns

```go
xlsx.New("test.xlsx", xlsx.SQLMode(xlsx.ModeBatch), xlsx.BatchSize(100), xlsx.TableName("test")).
  Header("id", "name").     // map columns
  Columns("key", "value").  // set key with value
  Build()
```

Output:
```SQL
INSERT INTO `test` (`id`, `name`, `key`) VALUES
('1', 'test01', 'value'),
('2', 'test02', 'value'),
('3', 'test03', 'value');
```

## License
© Wimi Yuan, 2023~time.Now <br>
Released under the [Apache License](./LICENSE).
