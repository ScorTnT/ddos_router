#!/bin/sh

# LAN 영역의 포트를 지정하는 설정입니다. 라우터의 하드웨어에 따라 다른 디바이스를 사용해도 됩니다.
uci -q del network.@device[0].ports
uci -q add_list network.@device[0].ports='eth1'
uci -q set network.lan.device='br-lan'

# WAN 영역의 포트를 지정하는 설정입니다. 위와 마찬가지로 다른 디바이스를 사용하셔도 됩니다.
uci -q set network.wan=interface
uci -q set network.wan.device='eth0'
uci -q set network.wan.proto='dhcp'

# 고정 아이피 세팅이 필요할 경우, 아래의 설정을 활용하세요.
# uci -q set network.wan=interface
# uci -q set network.wan.device='eth0'
# uci -q set network.wan.proto='static'
# uci -q set network.wan.ipaddr='220.66.59.116'
# uci -q set network.wan.netmask='255.255.255.0'
# uci -q set network.wan.gateway='220.66.59.2'
# uci -q set network.wan.dns='1.1.1.1 8.8.8.8'

# 아래는 와이파이 세팅입니다. 필요가 없는 경우, 언급하는 곳까지 전부 주석처리하세요.
[ -f /etc/config/wireless ] || wifi config

uci -q set wireless.radio0.country='KR'
uci -q set wireless.radio0.band='5g'
uci -q set wireless.radio0.channel='36'
uci -q set wireless.radio0.htmode='VHT80'
uci -q set wireless.radio0.disabled='0'

uci -q delete wireless.default_radio0 2>/dev/null
uci -q add wireless wifi-iface
uci -q rename wireless.@wifi-iface[-1]='default_radio0'
uci -q set wireless.default_radio0.device='radio0'
uci -q set wireless.default_radio0.network='lan'
uci -q set wireless.default_radio0.mode='ap'
uci -q set wireless.default_radio0.ssid='RPI-AP'
uci -q set wireless.default_radio0.encryption='psk2'
# 비밀번호는 본인의 상황에 맞게 변경하여 사용하세요.
uci -q set wireless.default_radio0.key='zoqtmxhs2025!'
# Management frame protection은 호환성 문제로 비활성화 (0). 필요하면 '1' 또는 '2'로 변경하세요.
uci -q set wireless.default_radio0.ieee80211w='0'
# 여기까지 주석처리하세요.

# Snort 패키지 기본 설정입니다. 필요 없는 경우 언급한 부분까지 주석처리하세요.
uci -q set snort.snort.enabled='1'
uci -q set snort.snort.manual='0'
uci -q set snort.snort.home_net='any'
uci -q set snort.snort.interface="$(uci get network.lan.device)"

uci commit snort
/etc/init.d/snort enable
/etc/init.d/snort start
# 여기까지 주석처리하세요.

uci commit wireless
uci commit network
/etc/init.d/network restart >/dev/null 2>&1 || true
wifi reload
exit 0
