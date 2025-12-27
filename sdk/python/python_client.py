import requests
import json

class Client:
    def __init__(self, addr: str):
        self.addr = addr.rstrip("/")

    # ---------- Single key ops ----------

    def SetKey(self, key, value):
        payload = {
            "key": key,
            "value": value
        }
        resp = requests.post(
            f"{self.addr}/set",
            json=payload
        )
        return resp.json()

    def GetKey(self, key):
        resp = requests.get(f"{self.addr}/get/{key}")
        return resp.json()

    def UpdateKey(self, key, value):
        payload = {
            "key": key,
            "value": value
        }
        resp = requests.put(
            f"{self.addr}/update",
            json=payload
        )
        return resp.json()

    def DeleteKey(self, key):
        resp = requests.delete(f"{self.addr}/delete/{key}")
        return resp.json()

    # ---------- Bulk ops ----------

    def GetAll(self):
        resp = requests.get(f"{self.addr}/all")
        return resp.json()

    def MSet(self, data: dict):
        """
        data = {
            "a": 10,
            "b": 20
        }
        """
        resp = requests.post(
            f"{self.addr}/mset",
            json=data
        )
        return resp.json()

    def MGet(self, keys: list[str]):
        """
        keys = ["a", "b", "c"]
        """
        resp = requests.post(
            f"{self.addr}/mget",
            json=keys
        )
        return resp.json()
