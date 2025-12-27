from sdk.python.python_client import Client

client = Client("http://localhost:9182")

client.SetKey("a", 10)
client.SetKey("b", 20)

print(client.GetKey("a"))
# {"value": 10}

print(client.MGet(["a", "b", "c"]))
# {"values": {"a": 10, "b": 20, "c": null}}

client.UpdateKey("a", 99)

print(client.GetAll())
# {"Data": {"a": 99, "b": 20}}

client.DeleteKey("b")
