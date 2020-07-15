# FileShare
For Sharing files using LAN

Requirement
- go get gopkg.in/cheggaaa/pb.v1
- go get github.com/grandcat/zeroconf

Then build using
go build main.go

To make file open to sharing run
- ./main -PATH "path_to_file"

This will out put a 5 digit code For eg- AB24D

to recieve the file run

- ./main -NICK CODE 

which here is AB24D so the actual command would be ./main -NICK AB24D
