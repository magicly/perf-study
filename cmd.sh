cc -O2 -o mc-pi-c mc-pi.c
c++ -O2 -o mc-pi-cpp mc-pi.cpp
time ./mc-pi-cpp
time ./mc-pi-c

go run monte_carlo.go
GOOS=windows go build monte_carlo.go