-- PREALPHA - NOT READY FOR USE YET --

# s3server

A Amazon S3-compatible server with interchangeable backends for development and testing purposes. Possibly production-ready in the future.

## Backends

- _File System_: Storing the buckets as directories and Objects as files
- _In Memory_: Stores everything in Main Memory state structures

## To get it properly working

To identify buckets S3 supports to methods: By path and by subdomain. That the s3server we need the ability to listen on a domain + subomdains. The easiest way to do this is dnsmasq

For setup on OSX follow [this tutorial](https://passingcuriosity.com/2013/dnsmasq-dev-osx/). 

The unit test will use test.dev as a domain. So dev should be bound to 127.0.0.1 by

    address=/dev/127.0.0.1



## Not supported features at the moment

