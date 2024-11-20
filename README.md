# att-fiber-gateway-info
## Description
A [golang](https://en.wikipedia.org/wiki/Go_(programming_language)) command line
tool to pull values from the pages of an
[AT&T Fiber](https://www.att.com/internet/fiber/) gateways.

## Supported hardware
* [BGW320-505 gateway](https://help.sonic.com/hc/en-us/articles/1500000066642-BGW320)
* BGW320-500 gateway, it has been reported to work
* BGW320-700 gateway, will likely work, but untested

## Supported firmware
I currently have version 6.30.5 on my
[BGW320-505 gateway](https://help.sonic.com/hc/en-us/articles/1500000066642-BGW320).
I have tested with previous versions back to 4.27.7, and expect them to work.

## Usage
```
att-fiber-gateway-info 1.0.11

Usage:
  -action string
        Action to perform (broadband-status, device-list, fiber-status, home-network-status,
                           ip-allocation, nat-check, nat-connections, nat-destinations,
                           nat-sources, nat-totals, reset-connection, reset-device,
                           reset-firewall, reset-ip, reset-wifi, restart-gateway,
                           system-information)
  -allmetrics bool (default: false)
        Return all metrics
  -cookiefile string (default: /var/tmp/att-fiber-gateway-info_cookies.gob)
        File to save session cookies
  -datadog bool (default: false)
        Send metrics to datadog
  -debug bool (default: false)
        Enable debug mode
  -filter string
        Filter to perform (icmp, ipv4, ipv6, tcp, udp)
  -fresh bool (default: false)
        Do not use existing cookies (Warning: If always used the gateway will
        run out of sessions.)
  -metrics bool (default: false)
        Return metrics instead of table data
  -no bool (default: false)
        Answer no to any questions
  -password string
        Gateway password
  -pretty bool (default: false)
        Enable pretty mode for nat-connections
  -statsdipport string
        Statsd ip port, like 127.0.0.1:8125
  -url string
        Gateway base URL
  -version bool (default: false)
        Show version
  -yes bool (default: false)
        Answer yes to any questions

```

## Compiling
To build for only system's platform:
```
go build
```

To binaries for all supported combinations of operating systems and architectures:
```
scripts/builds.sh
```

## Builds
See the `.go_builds` file for the list of supported combinations of operating
systems and architectures.

## Tests
```
# Run all tests
scripts/tests.sh

# Run all tests
scripts/tests.sh all

# Run tests for actions that do not require a login
scripts/tests.sh nologin

# Run tests for actions that are known to have metrics, and the -allmetrics flag
scripts/tests.sh metrics

# Run tests for actions that require login, but not the reset or restart actions
scripts/tests.sh login

# Run tests for the reset or restart actions
scripts/tests.sh reset
```

## Color
Colors for some of the output has been added, if the terminal supports it.

It can be disabled with `TERM=` or `NO_COLOR=`. `NO_COLOR` can be set to any
value other than empty to disable colors. This `NO_COLOR` any value behavior
seems to come from the shell, and isn't written into this code.

Examples:
```
TERM= ./att-fiber-gateway-info --help
NO_COLOR=1 ./att-fiber-gateway-info --help
NO_COLOR=false ./att-fiber-gateway-info --help # When set to anything, no color
```

## Configuration file
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

## Actions
See [ACTIONS.md](ACTIONS.md) in this git repository.

## Metrics
See [METRICS.md](METRICS.md) in this git repository.

## Story
  I just had [AT&T Fiber](https://www.att.com/internet/fiber/) installed. As
part of the installation process I received a
[AT&T](https://www.att.com/)
([Nokia](https://www.nokia.com/))
[BGW320-505 gateway](https://help.sonic.com/hc/en-us/articles/1500000066642-BGW320).
I immediately found and set up the
[passthrough mode](https://www.devonstephens.com/how-to-enable-ip-passthrough-on-att-bgw320-505/).
I already had my own [iptables](https://en.wikipedia.org/wiki/Iptables) based
router running [Fedora Linux](https://fedoraproject.org/). I wanted to
continue to use it. I wasn't interested in disabling my existing
[DHCP](https://en.wikipedia.org/wiki/Dynamic_Host_Configuration_Protocol)
servers and [DNS](https://en.wikipedia.org/wiki/Domain_Name_System) servers. I
also wasn't interested in changing the ip addresses of my whole network.

In my research one downside I found of
[passthrough mode](https://www.devonstephens.com/how-to-enable-ip-passthrough-on-att-bgw320-505/)
is that the
[BGW320-505 gateway](https://help.sonic.com/hc/en-us/articles/1500000066642-BGW320)
continues to do some connection tracking. It has a max of **8192** connections.
This can be tracked on the
[Diagnostics|NAT Table](https://192.168.1.254/cgi-bin/nattable.ha) page. It is
said that the most common way this becomes a problem is via
[bittorrent](https://en.wikipedia.org/wiki/BitTorrent) clients opening enough
connections to go over the limit.

The problem with the
[Diagnostics|NAT Table](https://192.168.1.254/cgi-bin/nattable.ha) page is it
is behind a [login](https://192.168.1.254/cgi-bin/login.ha) page. I wanted to
write a script to scrape the page. I then dug into the
[HTML](https://en.wikipedia.org/wiki/HTML) form and
[javascript](https://en.wikipedia.org/wiki/JavaScript). What I found is the
authors don't trust the network given you can access it via
[HTTP](https://en.wikipedia.org/wiki/HTTP) or
HTTPS](https://en.wikipedia.org/wiki/HTTPS) with
an invalid certificate. So they made the login process complicate.

There seems to a bug in the page design where the first time you go it will
claim your browser doesn't accept cookies, but it works on reload. My guess
is that this happens because it is trying to read the cookie that it hasn't
given you yet.

The login process goes something like this.

1. Load a page, get redirected to the
[login](https://192.168.1.254/cgi-bin/login.ha) page, and get told your browser
doesn't accept cookies.
2. Reload the page.
3. Receive the [nonce](https://en.wikipedia.org/wiki/Cryptographic_nonce)
hidden value in the page, and hashing functions written in
[javascript](https://en.wikipedia.org/wiki/JavaScript) that use
[md5](https://en.wikipedia.org/wiki/MD5).
[md5](https://en.wikipedia.org/wiki/MD5)(password+
[nonce](https://en.wikipedia.org/wiki/Cryptographic_nonce)).
4. Posting to the [login](https://192.168.1.254/cgi-bin/login.ha) page these
variables.
  - [nonce](https://en.wikipedia.org/wiki/Cryptographic_nonce), an all lower
case 64 character hex-decimal value. Looks like a
[sha256sum](https://en.wikipedia.org/wiki/SHA-2), like
`87428fc522803d31065e7bce3cf03fe475096631e5e07bbd7a0fde60c4cf25c7`.
  - password, not the actual password, but instead the same number of
characters as the password replaced with `*` characters
  - hashpassword, the results of the hashpwd function
  - Continue, the submit button that calls the hashing function and triggers
the post.

Once you are logged in your session is only tracked by your session cookie. It
is possible to generate too many sessions. I accomplished this during the
writing of this code, because at first it was generating a new session per run.

### Programming languages
I originally prototyped this written in
[javascript](https://en.wikipedia.org/wiki/JavaScript). Since I was trying to
reuse the [javascript](https://en.wikipedia.org/wiki/JavaScript) from the login
page. I started with [PhantomJS](https://github.com/ariya/phantomjs). It
worked, and was using [Firefox](https://www.mozilla.org/en-US/firefox/). This
meant [Firefox](https://www.mozilla.org/en-US/firefox/) needed to be installed,
and it needed to be able to run in the background. I switched to
[Puppeteer](https://github.com/puppeteer/puppeteer) once I realized
[PhantomJS](https://github.com/ariya/phantomjs) was not being maintained.
[Puppeteer](https://github.com/puppeteer/puppeteer) worked better, but suffered
the same problems as [PhantomJS](https://github.com/ariya/phantomjs).

I then realized I could just rewrite the
[javascript](https://en.wikipedia.org/wiki/JavaScript) code in
[golang](https://en.wikipedia.org/wiki/Go_(programming_language)). This works
well. It greatly simplifies the dependecies, CPU resources, memory resources,
testing, and more.

## Flaws in the login page
There are two flaws in the login process.

1. The security of the whole process is tied to the security of the session
cookie. Given that it goes across
[HTTP](https://en.wikipedia.org/wiki/HTTP)(no encryption) or
[HTTPS](https://en.wikipedia.org/wiki/HTTPS) with an invalid
certificate(easy to man in the middle), this is not good.
2. The hashing method is just a fig leaf, because of flaw one. In addition it
uses [md5](https://en.wikipedia.org/wiki/MD5), which is known to be insecure.
Modern solutions would be things like [sha256](https://en.wikipedia.org/wiki/SHA-2)
or [bcrypt](https://en.wikipedia.org/wiki/Bcrypt).

As far as I can tell there is no way to replace the gateway's
[SSL](https://en.wikipedia.org/wiki/Transport_Layer_Security) certificate.

I first ran across the first flaw over 15 years ago. I had gone to a
[Defcon](https://en.wikipedia.org/wiki/DEF_CON) talk about it,
and then discovered the same flaw in the real world.
