#!/bin/bash

# Initial name ID
name_id=1

# Loop over the range of addresses
for i in $(seq 1 10); do

    # Loop over the types
    for type in demo_25per.fastq demo_50per.fastq demo_75per.fastq demo.fastq; do

        # Loop for sending the POST request 10 times for each type
        for j in $(seq 1 10); do
            # Build the URL
            address="dgvkh-ven${i}.nrp-nautilus.io"

            # Send the POST request
            curl -X POST "http://${address}:8080/workflow" -H "Content-Type: application/json" -d '{
                "name": "workflow'"${name_id}"'",
                "type": "'"${type}"'",
                "received_at": "2023-07-23T12:34:56Z"
            }'
            sleep 0.2
            # Increment the name ID
            name_id=$((name_id + 1))
        done

    done

done
