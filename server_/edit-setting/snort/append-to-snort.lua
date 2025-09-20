ips = {
        mode = tap,
        rules = [[
                alert icmp 192.168.2.0/24 any -> any any (msg:"INLINE ICMP"; sid:1; rev:1;)
        ]]
}

--              alert tcp any any -> any any (msg:"INLINE TCP Test"; sid:1000002; rev:1;)
--              include /etc/snort/rules/local2.rules
--              include /etc/snort/rules/local.rules

daq = {
        module_dirs = { '/usr/lib/daq' },
        modules =
        {
                {
                        name = 'afpacket',
                        mode = 'inline'
                }
        }
}

