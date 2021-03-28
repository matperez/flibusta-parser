Парсер Флибусты
---------------

Флибуста раздает свою БД направо и налево по ссылке [https://flibusta.is/sql](https://flibusta.is/sql), но, если этого не знать, можно попробовать написать свой парсер.

ID книг растут монотонно с единицы до максимального значение, которое можно посмотреть на странице с последними поступлениями [https://flibusta.is/sql](https://flibusta.is/sql).

Отдельные книги доступны только авторизованным пользователям, так что придется добавить авторизацию.

## Сборка

Собирать можно через Makefile
```shell
make build
```

## Парсинг

Запуск через консоль

```shell
Usage: parser --db-user=STRING --db-password=STRING --flibusta-user=STRING --flibusta-password=STRING <command>

https://flibusta.is parser

Flags:
  -h, --help                          Show context-sensitive help.
      --db-server="localhost:3306"    Database server address and port
      --db-name="flibusta"            Database name
      --db-user=STRING                Database user name
      --db-password=STRING            Database user password
      --flibusta-user=STRING          Flibusta user name
      --flibusta-password=STRING      Flibusta user password

Commands:
  parse --db-user=STRING --db-password=STRING --flibusta-user=STRING --flibusta-password=STRING <from> <to>
    Run parsing.

Run "parser <command> --help" for more information on a command.

parser: error: missing flags: --db-user=STRING, --db-password=STRING, --flibusta-user=STRING, --flibusta-password=STRING

```
