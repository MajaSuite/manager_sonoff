# Manager instance for sonoff (wifi) devices

## Instead of preface

Manager is a driver module to serve requests and statuses from different hardware (sonoff in current case). It can receive
status notification from all sonoff hardware (diy and proprietary revisions). Proprietary hardware is encrypt exchange between 
device and cloud server, but also it encrypt status nofications via mDNS also (are you crazy, guys from sonoff?).

While development I was test integration using sonoff b1 lamp abd basic R2 relay, but I hope all kind of sonoff hardware will
communicate the same way.

All updates and commands send thru mqtt server. I assume different module (named hub) should control mqtt traffic and 
serve automation.

## How to sonoff (and this manager) works (in short)

Operation principle: Device send mDNS requests and put in TXT records sets of current device status. So manager receive 
updates and send it to connected mqtt server. To update device status manager receives mqtt push with field "cmd" inside 
json payload. Send request to device and after publish back update with device status.

This module use prefix sonoff/*device-id* in mqtt topic.

### Encryption  

I don't know encryption algorithm, but I hope it based on some device identification with float initialisation vector. 
Vector can change from time to time while communication. Addition salt is sequence number. Will try to reverse this a 
bit later, but not sure about sucess.

## Supported devices

B1 lamp (encrypted)
BASIC R2 relay

## Command line paraments

$ manager_sonoff -?
flag provided but not defined: -?
Usage of ./manager_sonoff:
  -clientid string
        client id for mqtt server (default "sonoff-1")
  -debug
        print debuging hex dumps
  -ip string
        sonoff server ip (default "192.168.1.1")
  -keepalive int
        keepalive timeout for mqtt server (default 30)
  -key string
        network key for registration (default "mypass")
  -login string
        login string for mqtt server
  -mqtt string
        mqtt server address (default "127.0.0.1:1883")
  -pass string
        password string for mqtt server
  -port string
        sonoff server port (default "8081")
  -reg
        to sonoff new device
  -sid string
        network name for registration (default "myhome")

## Device registration (-reg param)

Manager can adopt new device in factory default state. If device was previously linked you should turn on and off five 
times with 1 second delay. Signal light should start fast flashing. On computer you should connect to wireless network 
IDEAD-xxxxxx (without password) and run manager binary with follow parameters:

$ ./manager_sonoff -reg -sid MYNET -key MyNetPass -server ip.of.my.server -port xxxx

Actually *server* and *port* not used right now, but should defined in start (because must passed to device). Feel free to 
use your server address and port 8081

After registration complete your will device connect to your local wifi. Computer will be disconnected from IDEAD-xxxx 
network (and probably connect to you home network). If device registe successfully you should find it in mqtt.

## The commands from mqtt server to manage device

Command format: { "cmd": "xxxx", "data1": "aaa", "data2":"bbbb"}

List of supported commands:

+---------+-------------+-----------------------------------------------------------------------------+
| command | parameters  | description                                                                 |
+---------+-------------+-----------------------------------------------------------------------------+
| toggle  |             | toggle (swap) device power status                                           |
| power   | status      | Change power to *status* (on/off)                                           |
| onstart | state       | Set state while power on (on/off/stay). Stay means same as before power off |
| signal  |             | Get wifi signal status (negative value). Less is better.                    |
| pulse   | stat, width | Status and pulsation delay                                                  |
| setup   | sid, pass   | Change wifi credentials (network name and pass)                             |
| prepare |             | Prepare for OTA update (unlock ota update)                                  |
| update  | url, sha256 | Http url with firmware and firmware sha256 hash                             |
| status  |             | To fetch full device status                                                 |
+---------+-------------+-----------------------------------------------------------------------------+

## Known problems

 * I use my own mqtt client and server implementation. Now implementation is not so well, but works. Plan to completely rewite client
 * Encryption not supported
 * Need to verify most of sonoff device. Help in this case will be very appretiated.

## License and author

This project licensed under GNU GPLv3.

Author Eugene Chertikhin
