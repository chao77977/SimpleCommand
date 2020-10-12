# SimpleCommand
A Golang library to implement system command.

## Install
$ go get github.com/chao77977/simpleCommand

## Usage

### 1. Run the command and print the output on console.
Example,
```
cmd := simpleCommand.Newcommand("ping", "baidu.com", "-c 3")
r, o, err := cmd.RunWithOutput()
fmt.Println(cmd)
```
r: return value of command
o: strings of output
e: error
fmt.Println(cmd):
```
Exec the command: /usr/bin/ping baidu.com -c 3
Status          : completed
ExitCode        : 0
Output          : PING baidu.com (220.181.38.148) 56(84) bytes of data.
64 bytes from 220.181.38.148 (220.181.38.148): icmp_seq=1 ttl=48 time=44.4 ms
64 bytes from 220.181.38.148 (220.181.38.148): icmp_seq=2 ttl=48 time=10.9 ms
64 bytes from 220.181.38.148 (220.181.38.148): icmp_seq=3 ttl=48 time=17.9 ms

--- baidu.com ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 2003ms
rtt min/avg/max/mdev = 10.937/24.455/44.490/14.452 ms
```

Output on console
```
PING baidu.com (220.181.38.148) 56(84) bytes of data.
64 bytes from 220.181.38.148 (220.181.38.148): icmp_seq=1 ttl=48 time=5.00 ms
64 bytes from 220.181.38.148 (220.181.38.148): icmp_seq=2 ttl=48 time=5.04 ms
64 bytes from 220.181.38.148 (220.181.38.148): icmp_seq=3 ttl=48 time=5.01 ms

--- baidu.com ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 4688ms
rtt min/avg/max/mdev = 5.000/5.017/5.042/0.083 ms
```

### 2. Run the command and silent on console.
Example,
```
cmd := simpleCommand.Newcommand("ping", "baidu.com", "-c 3")
r, o, err := cmd.Run()
```
r: return value of command
o: strings of output
e: error

### 2. Run the command with timeout.
Example,
```
cmd := simpleCommand.Newcommand("ping", "baidu.com", "-c 3")
cmd.SetTimeoutWithSecond(1)
r, o, err := cmd.Run()
```
r: return value of command
o: strings of output
e: error
