# Configuration file

There is now a configuration file stored in the user's home directory as
`.att-fiber-gateway-info.yml` for Linux and MacOS. For Windows it is
`att-fiber-gateway-info.yml`. The file is automatically created if it does not
exist.

The default file permissions are read-write for the user only on Linux and MacOS.

The baseURL, with a default of `https://192.168.1.254`, is in the configuration
file. The command line argument overrides the configuration file. This value in
the configuration file is not optional.

The password is also in the configuration file. It is optional, and it is
commented out by default. It is recommended to set the password in the
configuration file for security. The command line argument overrides the
configuration file.
