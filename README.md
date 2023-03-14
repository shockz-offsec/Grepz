# Grepz
Grepz is a versatile command-line tool that offers various options for searching and highlighting specific terms or regex in a text file. With its customizable color options, users can highlight matches with a chosen background or foreground color. Additionally, Grepz allows users to view the context of the match by displaying the lines before and after it.

For reporting purposes, Grepz offers a special mode that highlights matches with a box, making it easy to spot them in a text file. Whether you need to quickly search for a term or carefully review a document, Grepz is a powerful and flexible tool that can help you achieve your goals.

## Installation
To install Grepz, run the following command:

```
go get github.com/shockz-offsec/grepz
```

```
git clone https://github.com/shockz-offsec/Grepz.git
cd Grepz
go build -o grepz grepz.go
```

## Download the compiled binary for Windows, Linux or MacOS

[Download the latest version]()

# 
Please, if you are going to use powershell on Windows 10 you must:

Activate Support VT (Virtual Terminal) / ANSI escape sequences globally by default, persistently.

```powershell
Set-ItemProperty HKCU:\Console VirtualTerminalLevel -Type DWORD 1
```


## Usage

```
grepz [-i] [-box|-b] [-up num] [-down num] <search_term> [<input_file>]
```

* `-i`: Performs a case-insensitive search.
* `-box` or `-b`: Highlights the matches inside a box.
* `-up num`: Displays a specified number of lines before the match.
* `-down num`: Displays a specified number of lines after the match.
* `<search_term>`: The term you want to search for. Regular expressions are accepted.
* `<input_file>`: The file in which you want to search. If this parameter is not provided, it is read from the standard input by pipping.
#
### Context funtionality

The context functionality allows grouping in the same context several occurrences that are within the range specified by the `-up` or `-down` parameters. This will make it easier to understand and analyze them together. In case the occurrences are not in range, they will be shown in different contexts, but always accompanied by their respective context lines for a better understanding.
#
### Colors

The background and foreground colors can also be changed by using the `-bg` and `-fg` flags followed by one of the following colors: `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, or `white`.

> The default background color is `magenta` and the foreground color is `white`.

> Box color is `red` by default.

## Examples

Search for the term `world` in the file `myfile.json`, using a case-insensitive search and highlight the match in `cyan`:

```sh
grepz -i -bg cyan "hello" myfile.json
```

[![](https://asciinema.org/a/Ic8TQu5ZsthyYPbFMJfwCi5p6.svg)](https://asciinema.org/a/Ic8TQu5ZsthyYPbFMJfwCi5p6)

Search for the term `Hello` in the file `myfile.json`, display 3 lines before and after the match:

```sh
grepz -up 3 -down 3 "Hello" myfile.json
```

[![](https://asciinema.org/a/2rjbKqRBlrHtBsYJTRdxGiVsJ.svg)](https://asciinema.org/a/2rjbKqRBlrHtBsYJTRdxGiVsJ)

Now the same example but highlighting the match with a box:

```sh
grepz -b -up 3 -down 3 "Hello" myfile.json
```

[![](https://asciinema.org/a/Olz4QnpNFofOfHha0di7FYfeb.svg)](https://asciinema.org/a/Olz4QnpNFofOfHha0di7FYfeb)

Search for the term `hello` in the `myfile.json` file, using a case-insensitive search, and highlight the match with a `yellow` box and `magenta` text:

```sh
grepz -i -b -bg yellow -fg magenta "hello" myfile.json
```

[![](https://asciinema.org/a/uuUia0T9FzvTEQpE5axiRHzRP.svg)](https://asciinema.org/a/uuUia0T9FzvTEQpE5axiRHzRP)

## ToDo
- [ ] Allow recursive searches in files

## Credits

[Shockz OffSec](https://github.com/shockz-offsec) & [Siriil](https://github.com/siriil)

## License

This tool is licensed under the Apache-2.0 License.