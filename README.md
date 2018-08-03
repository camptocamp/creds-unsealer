Creds-unsealer
==============

Generate your credentials files from secret stored in backends.

Usage
-----

```shell
Usage:
  creds-unsealer [OPTIONS]

Application Options:
  -V, --version                        Display version.
  -l, --loglevel=                      Set loglevel ('debug', 'info', 'warn', 'error', 'fatal', 'panic'). (default: info) [$BIVAC_LOG_LEVEL]
  -m, --manpage                        Output manpage.
  -b, --backend=                       Backend to use. (default: pass) [$CREDS_BACKEND]
  -p, --providers=                     Providers to use. (default: ovh, aws, openstack) [$CREDS_PROVIDERS]
      --output-key-prefix=             String to prepend to key of the secret
      --output-path-basedir=           Output path base directoty to use. (default: $HOME) [$OUTPUT_PATH_BASEDIR]

Pass backend options:
      --backend-pass-path=             Path to password-store. [$CREDS_BACKEND_PASS_PATH]

OVH Provider options:
      --provider-ovh-input-path=       OVH Provider input path (default: ovh)

AWS Provider options:
      --provider-aws-input-path=       AWS Provider input path (default: aws)

Openstack Provider options:
      --provider-openstack-input-path= Openstack Provider input path (default: openstack)

Help Options:
  -h, --help                           Show this help message
```

Backends
--------

The only supported backend for now is [pass](https://www.passwordstore.org/).

Providers
---------

### OVH

Default input path: `ovh`
Default output file: `~/.ovh.conf`
