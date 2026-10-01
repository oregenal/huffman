# Huz
Toy text archiver implemented using Huffman algorithm.  
## How to build
```console
go build hus.go
```
## How to use
**huz** have no special keys or whatever, it just automatically detect given file type, and perform corresponding
action.  
If it archive, **huz** will unpack it.
```console
huz archive.huz
```
Accordingly if it text file, **huz** will pack it.
```console
huz textfile.txt
```
In both cases the result will be saved as a separate file.
