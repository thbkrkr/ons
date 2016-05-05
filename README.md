# ons

```
Utility to manage OVH DNS zone.

Usage:
  ons COMMAND [arg...]

Available Commands:
  add         Plan to add a record
  apply       Changes DNS
  ls          List all DNS records of the zone
  plan        Show the execution plan
  rm          Plan to remove records matching a sub domain

Environment variables required:
  ONS_ZONE
  ONS_ENDPOINT, ONS_AK, ONS_AS, ONS_CK
```

## Setup

    mkdir dns

    # set credentials
    echo 'ONS_ZONE=bada.boum
    ONS_ENDPOINT=ovh-eu
    ONS_AK=LhLpR**********
    ONS_AS=yyYI2ZfCSC**********************
    ONS_CK=A2rkibPKri**********************' > dns/ons.env

    # bootstrap config
    echo '[]' > dns/ons.config.json

    # source .env
    eval $(cat dns/ons.env | sed "s:^:export :")

## Add a DNS record

    > ons add bim 1.2.3.4
    Refreshing DNS zone state prior to plan...

    + dns record: 1.2.3.4          bim.bada.boum

    Plan: 1 to add, 0 to remove.

    > ons apply
    Refreshing DNS state prior to apply...

    1.2.3.4          bim.bada.boum  added

    Apply: 1 added, 0 removed.

    > ons ls
    1.2.3.4               * bim.bada.boum