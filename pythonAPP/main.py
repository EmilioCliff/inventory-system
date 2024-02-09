import requests
url = "http://0.0.0.0:8080/products/"

header = {
    "Authorization": "Bearer v2.local.O9CyKcWYvVdGJm_6Sl8z1hRniGG03uBnJLg4_WxEQhfkHi-rW0fHAfq9YsR8ZknS80ItGAR9ecROdk0rzby6kNeJBk6BgyXoM6nYcFWzlscAH9wts_PlQ8yQiuMT3Dic-wGUPj9b7_NdqqY8hFn0MYZEemAl97sGPpnWYYh_g3zb-OMgcsNBsfsV4LqB8-hRlFzSqAZ7ll3Hi3ChnBYqWdulrrmEZgvYye69XfIzB-8skQfB6jYs85hH-4uBoUpNVtsNMHA7OiJcmwNlZAjLlg.bnVsbA"
}

response = requests.get(url=url, headers=header)

print(response.status_code)
print(response.text)