# 1. Enable IPv4 forwarding
sudo sysctl -w net.ipv4.ip_forward=1

# 2. Start/Up your tailscale daemon, advertising the Docker tailnet:
sudo tailscale up \
  --authkey tskey-… \
  --advertise-routes=10.46.0.0/16 \
  --accept-dns=true

# 3. Add NAT & FORWARD rules so containers on br-tailnet go out via tailscale0
sudo iptables -t nat -A POSTROUTING \
  -s 10.46.0.0/16 -o tailscale0 -j MASQUERADE
sudo iptables -A FORWARD \
  -i br-tailnet -o tailscale0 -j ACCEPT
sudo iptables -A FORWARD \
  -i tailscale0 -o br-tailnet -m state --state RELATED,ESTABLISHED -j ACCEPT
