# Create a header for the final CSV file
echo "Pod Name,Type,Iteration,RAM,CPUS,Duration" > durations.csv

echo "CSV header written to durations.csv"

# Define a function to handle the execution and data extraction for each pod
process_pod() {
    POD=$1
    FILETYPE=$2
    ITERATION=$3

    echo "Fetching RAM and CPUS environment variables from pod $POD..."
    RAM=$(kubectl exec $POD -- printenv RAM)
    CPUS=$(kubectl exec $POD -- printenv CPUS)

    echo "Executing workflow in $POD for file type $FILETYPE, Iteration $ITERATION..."
    kubectl exec $POD -- ./workflow/rna.sh $FILETYPE
    DURATION=$(kubectl exec $POD -- cat workflow/data/generated/result.txt)
    echo "$POD,$FILETYPE,$ITERATION,$RAM,$CPUS,$DURATION" >> "$POD"_"$FILETYPE"_"$ITERATION".output

    echo "Data saved to $POD"_"$FILETYPE"_"$ITERATION".output
}

# Export the function so it can be used by parallel
export -f process_pod

# File types to process
FILETYPES=("demo.fastq" "demo_25per.fastq" "demo_50per.fastq" "demo_75per.fastq")

# Get only "Running" backend pods and process them in parallel
for FILETYPE in "${FILETYPES[@]}"; do
    echo "Processing file type: $FILETYPE"
    kubectl get pods --field-selector=status.phase=Running | grep backend | awk '{print $1}' | \
    xargs -I {} -P 10 bash -c 'for i in {1..10}; do echo "Starting iteration $i for pod {} and file type $FILETYPE"; process_pod "$@" $i; echo "Finished iteration $i for pod {} and file type $FILETYPE"; done' _ {} $FILETYPE
done

# Aggregate results into a single CSV file
echo "Aggregating individual outputs into durations.csv..."
find . -name "*.output" -exec cat {} \; >> durations.csv

# Clean up temporary files
echo "Cleaning up temporary output files..."
rm *.output

# Display the results table
echo "Displaying final aggregated results:"
cat durations.csv
