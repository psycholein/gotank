go generate
GOOS=linux GOARCH=arm GOARM=6 go build
#ssh -l pi 192.168.0.226 "sudo killall gotank && cd ~/gotank && ./reset.sh"
scp -r gotank pi@192.168.0.208:~/gotank/
#scp ./config/*.yml pi@192.168.0.226:~/gotank/config/
#ssh -l pi 192.168.0.226 "cd ~/gotank && ./reset.sh"
#ssh -n -f -l pi 192.168.0.226 "sh -c 'cd ~/gotank; sudo nohup ./gotank > /dev/null 2>&1 &'"
