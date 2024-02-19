The fixtures in this directory are taken from:

https://github.com/isc-projects/bind9/tree/366b93f83565902aa98b3e146e0ead88e7502d87/tests/dns/testdata/master

We can validate input files by:

1. `docker run --rm -it --name bind9 -v $(pwd):/app debian:stable-slim`
2. `apt update && apt install -y bind9-utils`
3. `named-compilezone -o - example.com <master_file>`

## Notes

1. `master2.txt` includes a bad record line: `a in ns`
2. `master3.txt` has 3 lines without owner fields (i.e. no leading `@` or domain)
3. `master4.txt` does not define a TTL, so it falls back to the SOA minttl
4. `master5.txt` has the "any" `QCLASS` in place of a valid `CLASS`
5. `master6.txt` has no `ns` records, so Bind would refuse to load it
6. `master7.txt` has incomplete values for the DNSKEY records (missing hashes)
7. `master17.txt` has out-of-zone data on lines 3 and 12. It also has an
inherited owner (`b.test`) immediately after an `$ORIGIN` directive (`sub.test`).
Bind will refuse to compile this file. We end up ignoring the directive.
