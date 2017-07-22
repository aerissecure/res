# res
Parallel hostname resolver with flexible output

# Examples

Normal usage:

~~~
> $ res yahoo.com 8.8.8.8 google.com
yahoo.com
   98.139.180.149
   206.190.36.45
   98.138.253.109
   2001:4998:58:c02::a9
   2001:4998:c:a06::2:4008
   2001:4998:44:204::a7

8.8.8.8
   google-public-dns-a.google.com.

google.com
   172.217.9.78
   2607:f8b0:4009:80e::200e
~~~

Filter for only IPv4 addresses:

~~~
> $ res -4 yahoo.com 8.8.8.8 google.com
yahoo.com
   98.139.180.149
   206.190.36.45
   98.138.253.109

8.8.8.8
   google-public-dns-a.google.com.

google.com
   172.217.9.78

~~~

Perform resolutions for lookup responses (recursive)

~~~
> $ res -r yahoo.com                                                                                [±master ●]
yahoo.com
   98.139.180.149
      ir1.fp.vip.bf1.yahoo.com.
   206.190.36.45
      ir1.fp.vip.gq1.yahoo.com.
   98.138.253.109
      ir1.fp.vip.ne1.yahoo.com.
   2001:4998:58:c02::a9
      ir1.fp.vip.bf1.yahoo.com.
   2001:4998:c:a06::2:4008
      ir1.fp.vip.gq1.yahoo.com.
   2001:4998:44:204::a7
      ir1.fp.vip.ne1.yahoo.com.
~~~

Res help

~~~
Usage of res:
  -4  ipv4 responses only
  -6  ipv6 responses only
  -j  output json
  -jp
      output json (pretty)
  -r  recursive (resolve lookups)
~~~