import json
import requests


def submit_workflow(ven_url, workflow_name,workflow_type):
    # Prepare the workflow to submit
    submit_workflow = {
        "Name": workflow_name,
        "Type": workflow_type,
    }
    body = json.dumps(submit_workflow)

    # Send the POST request
    try:
        response = requests.post(ven_url + "/workflow", headers={"Content-Type": "application/json"}, data=body)
    except requests.RequestException as e:
        raise Exception(f"HTTP request failed: {e}")
    print(response)
    print(json.loads(response.text)["id"])


    # Check the response status
    if response.status_code == 451:  # StatusUnavailableForLegalReasons
        raise Exception(f"unexpected status code: {response.status_code}")

    if response.status_code != 201:  # StatusCreated
        raise Exception(f"unexpected status code: {response.status_code}")


submit_workflow("https://dgvkh-ven3.nrp-nautilus.io", "aaaa","demo.fastq")