go build -o nineveh.so -x -buildmode=c-shared main.go
mv nineveh.so nineveh.h /home/ajmollohan/alucard


#go build -o nineveh_one_step.so -buildmode=c-shared -v main.go
#mv nineveh_one_step.so nineveh_one_step.h /home/ajmollohan/alucard
