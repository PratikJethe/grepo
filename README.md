# Grepo
grepo is an implementation of the Unix command grep in Go. grep is a tool used for searching files or input streams for lines that match a specified pattern.

## Features

- search in provided file
- search in directory
- performs case sensitive,case insensitive, exact matches
- stores output in given file path
- output match count, lines after and before match

## How to use
- clone this repository [git clone https://github.com/PratikJethe/grepo.git]
- run go build in root of the project
- execute binary with appropriate flags [./grepo.exe -f filename.txt -s search_word]


```sh
npm install --production
NODE_ENV=production node app
```

## Flags

grepo supports various flags. multiple flags can be combined to generate desired output

| flag | desrcription | datatype
| ------ | ------ | ---- |
| -f | accepts file path where search is to be done | string
| -s | accepts search query| string
| -i | performs case insensitive search| bool
| -e | performs exact matching search | bool
| -input | accepts word from user | bool
| -o | accepts output file path to store output | string
| -dir | accepts path to directory where search is to be performed | string
| -c | show only count of matches | bool
| -a | display lines after match | bool
| -b | display lines before match | bool


## Examples 
###### Note: Test data used for examples is taken from test folder in repository 
#
#
### 1. Simple search in given file
 
```sh
$ ./grepo.exe -f test/test_data.txt -s test
Match in file: test/test_data.txt  line  1:0 "test line for one occurence in single line"
Match in file: test/test_data.txt  line  2:0 "test for multiple occurences in single line test test"
Match in file: test/test_data.txt  line  2:44 "test for multiple occurences in single line test test"
Match in file: test/test_data.txt  line  2:49 "test for multiple occurences in single line test test"
```

### 2. Case insensitive search in given file
 
```sh
$ ./grepo.exe -f test/test_data.txt -s test -i
Match in file: test/test_data.txt  line  1:0 "test line for one occurence in single line"
Match in file: test/test_data.txt  line  2:0 "test for multiple occurences in single line test test"
Match in file: test/test_data.txt  line  2:44 "test for multiple occurences in single line test test"
Match in file: test/test_data.txt  line  2:49 "test for multiple occurences in single line test test"
Match in file: test/test_data.txt  line  4:0 "Test for case-insensitive "
Match in file: test/test_data.txt  line  5:0 "Testing for excat match"
```
### 3. Exact matching search in given file
 
```sh
$ ./grepo.exe -f test/test_data.txt -s Testing  -e
Match in file: test/test_data.txt  line  5:0 "Testing for excat match"
```
### 4. -i -e flags combined
 
```sh
$ ./grepo.exe -f test/test_data.txt -s testing  -e -i
Match in file: test/test_data.txt  line  5:0 "Testing for excat match"
```
### 5. Accept input from user
 
```sh
$ ./grepo.exe -input -s test
lorem testing Testing test
Match found: testing
Match found: test
```

### 6. Accept input from user with -i and -e flag
 
```sh
$ ./grepo.exe -input -s test -i -e
lorem testing Testing test Test
Match found: test
Match found: Test
```

### 7. Store output in given filepath
 
```sh
$ ./grepo.exe -f test/test_data.txt -s test -o output.txt
output stored into output.txt
$ cat output.txt 
Match in file: test/test_data.txt  line  1:0 "test line for one occurence in single line"
Match in file: test/test_data.txt  line  2:0 "test for multiple occurences in single line test test"
Match in file: test/test_data.txt  line  2:44 "test for multiple occurences in single line test test"
Match in file: test/test_data.txt  line  2:49 "test for multiple occurences in single line test test"
```

### 8. Search in all text files of a given directory
 
```sh
$ ./grepo.exe -dir test -s  test
Match in file: test\directory_one\test_data_dir_1.txt  line  1:0 "test data in directory one"
Match in file: test\directory_two\test_data_dir_2.txt  line  1:0 "test data in directory one"        
Match in file: test\test_data.txt  line  1:0 "test line for one occurence in single line"
Match in file: test\test_data.txt  line  2:0 "test for multiple occurences in single line test test" 
Match in file: test\test_data.txt  line  2:44 "test for multiple occurences in single line test test"
Match in file: test\test_data.txt  line  2:49 "test for multiple occurences in single line test test"
```

### 9. Show only count of matches
 
```sh
$ ./grepo.exe -dir test -s  test -c
Number of matches : 6
```

### 10. Show lines before match
 
```sh
$ ./grepo.exe -s test  -f test/before_after.txt  -b
line number 1 above match 
line number 2 above match
line number 3 above match
```
### 11. Show lines before match
 
```sh
$ ./grepo.exe -s test  -f test/before_after.txt   -b
line number 1 before match 
line number 2 before match
line number 3 before match
```
### 12. Show lines after match
 
```sh
$ ./grepo.exe -s test  -f test/before_after.txt   -a
line number 1 after match 
line number 2 after match
line number 3 after match
```
### 13. Show lines after and before match
 
```sh
$ ./grepo.exe -s test  -f test/before_after.txt   -a -b
line number 1 before match 
line number 2 before match
line number 3 before match
line number 1 after match
line number 2 after match
line number 3 after match
```
