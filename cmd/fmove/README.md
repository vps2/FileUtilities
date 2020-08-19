# FMove

**FMove** - tracks the appearance of files in one folder and moves them to another.

```sh
fmove.exe -h
Usage of fmove.exe:
  -s, --src-dir string      the folder where new files are tracked
  -d, --dst-dir string      the folder where new files will be moved from the source folder
  -t, --timeout duration    the timeout between polls of the source directory (default 1m0s)
```

### Usage example:

```sh
fmove.exe -s c:\temp\1 -d c:\temp\2 -i 15s
```