Creds-unsealer
==============

Generate your credentials files from secret stored in backends.

Usage
-----

```shell


Usage:
  creds-unsealer [OPTIONS]

Application Options:
  -V, --version                  Display version.
  -l, --loglevel=                Set loglevel ('debug', 'info', 'warn', 'error', 'fatal', 'panic').
                                 (default: info) [$BIVAC_LOG_LEVEL]
  -m, --manpage                  Output manpage.
  -b, --backend=                 Backend to use. (default: pass) [$CREDS_BACKEND]
  -p, --providers=               Providers to use. (default: ovh) [$CREDS_PROVIDERS]
      --output-key-prefix=       String to prepend to key of the secret

Pass backend options:
      --backend-pass-path=       Path to password-store. [$CREDS_BACKEND_PASS_PATH]

OVH Provider options:
      --provider-ovh-input-path= OVH Provider input path (default: ovh)

Help Options:
  -h, --help                     Show this help message
```

Backends
--------

The only supported backend for now is [pass](https://www.passwordstore.org/).

Providers
---------

### OVH

Default input path: `ovh`
DEfault output file: `~/.ovh.conf`
