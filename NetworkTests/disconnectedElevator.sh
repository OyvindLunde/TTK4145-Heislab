elevatorPort1=11111
elevatorPort2=11112

sudo iptables -A INPUT -p tcp --dport $elevatorPort1 -j ACCEPT
sudo iptables -A INPUT -p tcp --sport $elevatorPort1 -j ACCEPT

sudo iptables -A INPUT -p tcp --dport $elevatorPort2 -j ACCEPT
sudo iptables -A INPUT -p tcp --sport $elevatorPort2 -j ACCEPT

sudo iptables -A INPUT -j DROP
