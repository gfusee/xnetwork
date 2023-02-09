import requests
import time

def wait_for_vm_query_ready():
    url = "http://testnet:7950/vm-values/query"
    payload1 = '{"args":[],"caller":"erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqplllst77y4l","funcName":"getQueueRegisterNonceAndRewardAddress","scAddress":"erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqllls0lczs7"}'
    payload2 = '{"args":[],"funcName":"getAllContractAddresses","scAddress":"erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqylllslmq6y6"}'
    headers = {
        'Content-Type': 'application/json'
    }

    while True:
        print("Querying the proxy...")
        response1 = requests.request("POST", url, headers=headers, data=payload1)
        response2 = requests.request("POST", url, headers=headers, data=payload2)
        if response1.status_code == 200 and response2.status_code == 200:
            break
        print("VM not ready for queries = Will retry in 1 second...")
        time.sleep(1)


wait_for_vm_query_ready()
