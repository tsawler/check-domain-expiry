# check-domain-expiry

A simple nagios module, written in go, to test for domain expiration.

This plugin is meant to be used with [Nagios](https://www.nagios.org/).


# Build
Compile for Linux (e.g. Digital Ocean Ubuntu 16.04): 

~~~
env GOOS=linux GOARCH=amd64 go build -o check_domain_expiration main.go
~~~

# Usage

Run the command from cli as follows:

~~~
check_domain_expiration -host <domainname.com>
~~~

## Integration with Nagios 4

Add this to `/usr/local/nagios/objects/commands.cfg`:

~~~
define command {
   command_name    check_domain_expiration
   command_line    /usr/local/nagios/libexec/check_domain_expiration -host $ARG1$
}
~~~


In individual files in `/usr/local/nagios/etc/servers`:

~~~
define service{
        use                     generic-service
        host_name               www.somesite.com
        service_description     Check Domain Expiration
        check_command           check_domain_expiration!domain.com
}

~~~