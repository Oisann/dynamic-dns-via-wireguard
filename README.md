# Dynamic DNS via Wireguard

Updates a Cloudflare DNS record based on the remote endpoint of a wireguard connection.

## Config

```yaml
records:
  -
    email: <cloudflare-dashboard-email>
    key: <wireguard-peer-public-key>
    name: <subdomain>.<domain>
    record: <cloudflare-dns-record-id>
    token: <cloudflare-api-token>
    ttl: <dns-ttl>
    zone: <cloudflare-zone-id>
    proxied: <true-or-false>
settings:
  interval: <seconds>
``` 

TTL set to `1` will set it to automatic.
The DNS Record ID can be hard to get, I found it by [curl'ing the DNS Record Details via the API9](https://api.cloudflare.com/#dns-records-for-a-zone-dns-record-details).
It should probably be part of this product, so feel free to make a pull request!
The interval is how often it will check the tunnel endpoints for a change.

## Flags

* --config \<config-file\>
