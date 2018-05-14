package main

// Raspberry Pi Server Constants
var (
    // PiTimeout Timeout in Seconds
    PiTimeout = 3

    // PiPollPeriod poll period in milliseconds
    PiPollPeriod = 1000

    // PiIPAddress the ipv4 address of the Raspberry Pi
    PiIPAddress = "169.254.56.46"

    // PiPort the port of the pi server
    PiPort = ":1337"
)

// Web Server Constants
var (
    // WebTimeout the timeout value in seconds of the webserver
    WebTimeout = 3

    // WebIPAddress the IPv4 address of the web server
    WebIPAddress = "169.254.159.115"

    // WebPort the port number of the web server
    WebPort = ":9001"
)
