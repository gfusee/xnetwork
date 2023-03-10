import requests
import time

def wait_for_vm_query_ready():
    query_url = "http://testnet:7950/vm-values/query"
    query_payload1 = '{"args":[],"caller":"erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqplllst77y4l","funcName":"getQueueRegisterNonceAndRewardAddress","scAddress":"erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqllls0lczs7"}'
    query_payload2 = '{"scAddress":"erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqylllslmq6y6","funcName":"getAllContractAddresses","args":[]}'

    economics_url = "http://testnet:7950/network/economics"

    headers = {
        'Content-Type': 'application/json'
    }

    while True:
        query_response1 = requests.request("POST", query_url, headers=headers, data=query_payload1)
        query_response2 = requests.request("POST", query_url, headers=headers, data=query_payload2)
        economics_response = requests.request("GET", economics_url, headers=headers)
        if query_response1.status_code == 200 and query_response2.status_code == 200 and economics_response.status_code == 200:
            break
        time.sleep(1)


wait_for_vm_query_ready()
