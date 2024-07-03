from urllib.parse import urlencode
from urllib.request import Request, urlopen

url = "http://127.0.0.1:8088/debug/test"  # Set destination URL here
post_fields = {"zone": "test", "command": "restart"}  # Set POST fields here

request = Request(url, urlencode(post_fields).encode())
with urlopen(request) as response:
    json = response.read().decode()
print(json)
