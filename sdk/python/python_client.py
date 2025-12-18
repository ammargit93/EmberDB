import requests 

class Client:
    def __init__(self, addr):
        self.addr = addr

    def SetKey(self, key, value):
        resp = requests.post(self.addr+"/set",  data='{"key": "%s", "value": %s}' % (key, value),headers={"Content-Type":"application/json"})     
        return resp.json()
    