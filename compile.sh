#todo: better automation, I manually change this every time my circumstances change
go build -o nineveh.so -x -buildmode=c-shared main.go
mv nineveh.so nineveh.h /home/ajmollohan/alucard

