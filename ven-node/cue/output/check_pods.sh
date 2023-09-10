# Get all backend pods and execute the command, redirecting output to separate files
kubectl get pods | grep backend | awk '{print $1}' | \
xargs -I {} -P 10 sh -c 'kubectl exec  {} -- ./workflow/rna.sh demo_25per.fastq > {}.output 2>&1'
