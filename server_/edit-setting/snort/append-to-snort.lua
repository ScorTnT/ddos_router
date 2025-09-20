ids = {
        mode = tap,
        rules = [[
                alert icmp 192.168.3.225/32 any -> any any (msg:"ICMP packet from attacker"; flags:S; sid:1; rev:1;)
                alert tcp 192.168.3.225/32 any -> any any (msg:"SYN packet from attacker to any port"; flags:S; sid:2; rev:1;)
                alert tcp 192.168.3.225/32 any -> any 80 (msg:"SYN packet from attacker"; flags:S; sid:3; rev:1;)
        ]]
}

daq = {
        module_dirs = { '/usr/lib/daq' },
        modules =
        {
                {
                        name = 'afpacket',
                        mode = 'passive'
                }
        }
}
