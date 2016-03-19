# res
Parallel hostname resolver with flexible output

# Examples

Normal usage:

~~~
> $ res yahoo.com google.com bing.com
yahoo.com
   98.139.183.24
   98.138.253.109
   206.190.36.45
   2001:4998:c:a06::2:4008
   2001:4998:44:204::a7
   2001:4998:58:c02::a9

bing.com
   204.79.197.200

google.com
   216.58.193.174
   2607:f8b0:400a:800::200e
~~~

Filter for only IPv4 addresses:

~~~
> $ res -4 yahoo.com google.com bing.com
google.com
   216.58.216.142

bing.com
   204.79.197.200

yahoo.com
   98.138.253.109
   98.139.183.24
   206.190.36.45
~~~

Output to a single line with space separated addresses:

~~~
> $ res -l -4 yahoo.com google.com bing.com
204.79.197.200 98.138.253.109 98.139.183.24 206.190.36.45 216.58.216.142
~~~
