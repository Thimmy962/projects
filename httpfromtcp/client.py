import requests

for i in range(2000): 
	res = requests.get("http://localhost:42069")
	print(res.text, end="")
