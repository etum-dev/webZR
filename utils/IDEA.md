You can make a config file of scope and specify how this shitty tool fuzzes it.
fuzzlevel 1-3, where 1 is most gentle.
It also has support for comments.
for example:

$ cat config
// These are for the XYZ customer
[fuzzlevel:1]
[header:i-am-not-a-badguy]
somedomain.com
someotherdomain.com 
// I hate these hoes
[fuzzlevel:3]
[header:Mozilla/4.05 (Macintosh; I; PPC, Nav)]
microsoft.com
google.com
tiktok.com