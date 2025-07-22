import requests

def call_connect_api(command: str) -> str:
    url = "http://localhost:8001/connect-db"
    try:
        response = requests.post(url, data=command, headers={"Content-Type": "text/plain"})
        return response.text
    except requests.exceptions.RequestException as e:
        print("Failed to call connect API:", e)
        return "Error: " + str(e)

output = call_connect_api("GET b")
print(output)
