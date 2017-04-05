BuffaLab
========

`buffalabd` is a small network server, written in Go, that
provides:

  - tFTP / PXE configuration files
  - HTTP installation media
  - A REST API for managing a lab

What is it for?

Imagine you have a small amount of hardware to use as a home lab,
and you want to automate the boring bits of installing operating
systems, formatting disks, installing packages, and configuration.

You could use something like Puppet or Ansible, sure, but you
still have to get something _up and running_ to leverage those
tools, and they require some supporting infrastructure.

Instead, the route buffalabd takes is to use the Pre-eXecution
Environment support in modern baseboard management agents to boot
the OS installation over the network, by way of DHCP + PXE + HTTP.
On top of this magic, `buffalabd` also provides an API for
managing PXE boot files, so that you have a completely hands-off
lab experience &mdash; no iLO or iDRAC web UIs in sight.


Configuration
-------------

`buffalabd` reads its configuration at startup, from
`buffalab.yml` in the current working directory.  You can supply a
different file by settings its path in the `$BUFFALAB_CONFIG`
environment variable.

Here is an example configuration:

```
---
# Where should the files for tFTP and HTTP be served from?
root: /srv/lab

# What port should the HTTP image server / REST API listen on?
listen: ':8085'

# What configuration flavors are supported?
flavors:
  - vanilla
  - chocolate

# What machines comprise your lab?
machines:
  - name:     lab1
    role:     leader
    mac:      aa-bb-cc-11-22-33
    type:     idrac
    ip:       10.0.0.12
    username: root
    password: calvin

  - name:     lab2
    role:     drone
    mac:      aa-bb-cc-11-22-34
    type:     idrac
    ip:       10.0.0.12
    username: root
    password: calvin
```


API Reference
-------------

Here are the endpoints you can hit on the API:

### GET /api/status

Retrieves the current status of the lab.

```
curl -q http://127.0.0.1:8085/api/status | jq -r .
{
  "flavor": "vanilla",
  "nodes": {
    "lab1": "ok",
    "lab2": "fail"
  }
}
```

### POST /api/install

Reconfigures the tFTP boot configuration for all machines in the
lab, and then power-cycles all of the machines to start the
installation process.

Provide the flavor name as the body of the POST data:

```
curl -X POST http://127.0.0.1:8085/api/install -d 'vanilla'
```

Returns an HTTP 204 on success, with no body.

### POST /api/:machine/fail

Indicates that the named `:machine` has failed in its
configuration steps.  For example:

```
curl -X POST http://127.0.0.1:8085/api/lab1/fail
```

Returns an HTTP 204 on success, with no body.

### POST /api/:machine/ok

Indicates that the named `:machine` has succeeded in its
configuration steps.  For example:

```
curl -X POST http://127.0.0.1:8085/api/lab2/fail
```

Returns an HTTP 204 on success, with no body.


Installation
------------

You'll need a recent vintage `Go` toolchain.  Then, all you have
to do is `make`.  Stay tuned for more information on the full
Raspberry Pi experience...


It's Still Evolving
-------------------

There's more features I want to add to this thing; so it's far
from feature-complete.  If you think this is something you'd like
to use, and want to help out, raise a [Github Issue][issues] and
we can discuss the implementation before you go off and submit a
Pull Request.




[issues]: https://github.com/jhunt/buffalabd
