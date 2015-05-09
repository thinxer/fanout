# fanout

Simple Go library for broadcasting a value to many receivers.

# Usage

    import "github.com/thinxer/fanout"
    
    var fan fanout.Fan
    
    go func() {
        encoder := json.NewEncoder(os.Stdout)
        err := fan.Receive(10, func(v interface{}) error {
            return encoder.Encode(v)
        })
        log.Println("Encoder error:", err)
    }()
    
    // send anything
    fan.Send(value)
  
  
