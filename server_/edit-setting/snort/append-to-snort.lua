ids = {
        mode = tap,
        rules = [[
                alert icmp !192.168.2.0/24 any -> any any (msg:"non-LAN source"; flags:S; sid:1; rev:1;)
                alert tcp !192.168.2.0/24 any -> any any (msg:"non-LAN source to any port"; flags:S; sid:2; rev:1;)
                alert tcp !192.168.2.0/24 any -> any 80 (msg:"non-LAN source to port 80"; flags:S; sid:3; rev:1;)
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



