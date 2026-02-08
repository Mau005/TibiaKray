# TibiaKray


WireGuard:

sudo cat /etc/wireguard/wg0.conf

sudo wg-quick down wg0
sudo wg-quick up wg0

sudo wg


sudo chmod 600 /etc/wireguard/wg0.conf



[Interface]
Address = 10.100.0.1/24
ListenPort = 51820
PrivateKey = privatekey

PostUp = sysctl -w net.ipv4.ip_forward=1; iptables -A FORWARD -i wg0 -j ACCEPT; iptables -A FORWARD -o wg0 -j ACCEPT; iptables -t nat -A PREROUTING -p tcp --dport 80 -j DNAT --to-destination 10.100.0.2:80; iptables -t nat -A PREROUTING -p tcp --dport 443 -j DNAT --to-destination 10.100.0.2:443
PostDown = iptables -D FORWARD -i wg0 -j ACCEPT; iptables -D FORWARD -o wg0 -j ACCEPT; iptables -t nat -D PREROUTING -p tcp --dport 80 -j DNAT --to-destination 10.100.0.2:80; iptables -t nat -D PREROUTING -p tcp --dport 443 -j DNAT --to-destination 10.100.0.2:443

[Peer]
PublicKey = publickey
AllowedIPs = 10.100.0.2/32

Arranque automatico
sudo systemctl enable wg-quick@wg0



Regla:
sudo iptables -t nat -A POSTROUTING -o wg0 -p tcp -d 10.100.0.2 -m multiport --dports 80,443 -j MASQUERADE


Cliente Internet
   ↓
VPS (eth0:80/443)
   ↓ DNAT
WireGuard (wg0)
   ↓
Casa (10.100.0.2:80/443)
   ↑ respuesta
VPS (MASQUERADE wg0)
   ↑
Cliente Internet