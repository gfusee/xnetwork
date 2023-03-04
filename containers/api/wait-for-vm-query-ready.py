import requests
import time

def wait_for_vm_query_ready():
    url = "http://testnet:7950/vm-values/query"
    payload = '{"args":[],"caller":"erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqplllst77y4l","funcName":"getQueueRegisterNonceAndRewardAddress","scAddress":"erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqllls0lczs7"}'
    headers = {
        'Content-Type': 'application/json'
    }

    while True:
        response = requests.request("POST", url, headers=headers, data=payload)
        if response.status_code == 200:
            break
        time.sleep(1)


wait_for_vm_query_ready()
