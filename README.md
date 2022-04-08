# [Ludum Dare 50](https://ldjam.com/events/ludum-dare/50/kernel-panic)

# Goal:
You are the administrator of this server.                                                                             
Some hackers plant viruses.                                                                                           
You have to stop them to avoid a kernel panic!

# Run it
https://kernel.aligator.dev/ssh/host/aligator.dev
username & password: 'admin'

or

ssh -p 2223 aligator.dev
username & password: 'admin'

Or using any other ssh client.

(Note any ssh key and any user/password will work - there is no check)

## Locally
`go run ./cmd/local`  
will start a local version (installed Go is required)

# Scrolling:                                                                                                            
In some terminals the scroll wheel just works.                                                                        
In others, just use the "PageUp" and "PageDown" or "ctrl+u" and "ctrl+d"                                              
                                                                                                                      
# Quit game:                                                                                                            
"ctrl+c" or type "exit"                                                                                               
                                                                                                                      
# Available commands:                                                                                                   
• ls {Path}    (list directory)                                                                                       
• mkdir {Path} (create folder)                                                                                        
• cd {Path}    (change directory)                                                                                     
• rm {Path}    (delete file)                                                                                          
• ps           (list processes)                                                                                       
• kill {PID}   (kill process)
