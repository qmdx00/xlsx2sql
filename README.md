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

### Single insert statements

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

### Batch insert statements

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
  Mapped("idx", 0).                       // get id from xlsx column 0
  Mapped("name", 1).                      // get name from xlsx column 1
  Columns("key", "value").                // set key with value
  Valuer("idx", func(v string) string {   // custom value processor
    id, _ := strconv.Atoi(v)
    return fmt.Sprintf("%d", id - 1)
  }).Build()
```

Output:
```SQL
INSERT INTO `test` (`idx`, `name`, `key`) VALUES
('0', 'test01', 'value'),
('1', 'test02', 'value'),
('2', 'test03', 'value');
```

### Mappping columns with Header

```go
xlsx.New("test.xlsx", xlsx.SQLMode(xlsx.ModeBatch), xlsx.BatchSize(100), xlsx.TableName("test")).
  Header("id", "name").     // mapping columns
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
Released under the [MIT License](./LICENSE).
