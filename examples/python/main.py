from sdk.python.python_client import Client

client = Client("http://localhost:9182")

# resp = client.SetKey("bodycount", 100)

val = client.GetKey("bodycount")
print(val)