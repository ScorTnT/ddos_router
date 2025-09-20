ips = {
        mode = tap,
        rules = [[
                alert icmp 192.168.3.225/32 any -> any any (msg:"INLINE ICMP"; sid:1; rev:1;)
                alert tcp 192.168.3.225/32 any -> any any (msg:"INLINE TCP"; sid:2; rev:1;)
        ]]
}

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

