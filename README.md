# gomorphy
pymorphy3 port for golang

Частичный порт на язык Go библиотеки python https://github.com/no-plagiarism/pymorphy3

Реализована поддержка только русского языка (украинский в планах на будущее).

Частичный форк https://github.com/jus1d/gomorphy . За основу взят код работы со словарем DAWG.

Работа со словарем opencorpora в формате DAWG (словарь от 22.05.2026).

(Скопируйте каталог словаря opencorpora из пакета gomorphy в папку со своим проектом.) 

Реализованы все эвристические анализаторы.

Реализован метод Parse(). Для удобства добавлен метод IsParsed()

Реализован метод Inflect(). Для удобства добавлен метод InflectVar() с вариадическими аргументами.

## License

The **Go source code** is licensed under the [MIT License](LICENSE).

The **embedded dictionary data** (`opencorpora/`) is derived from [OpenCorpora](http://opencorpora.org/)
and is licensed under [CC BY-SA 4.0](opencorpora/LICENSE). If you distribute a binary that embeds
this data, you must comply with the CC BY-SA 4.0 terms (attribution + ShareAlike).
