#!/bin/bash
source /miniconda/etc/profile.d/conda.sh && conda activate rnaseq

set -e
set -o pipefail

SECONDS=0

WORKFLOW_DIR="./workflow"
DATA_DIR="${WORKFLOW_DIR}/data"
# /app/workflow/data/generated
GENERATED_DIR="${DATA_DIR}/generated"
FASTQ_FILE="demo.fastq"
TRIMMED_FILE="demo_trimmed.fastq"
TRIMMED_BAM="demo_trimmed.bam"
FEATURE_COUNTS="demo_featurecounts.txt"
DURATION_FILE="result.txt"
LOG_FILE="logs.txt"



# Define error trap
function handle_error() {
    echo "Error occurred in the script at line $LINENO during execution of command: '$BASH_COMMAND'. Exiting with status $?. Output: 'failed'" > ${GENERATED_DIR}/$DURATION_FILE
    echo "Error occurred in the script at line $LINENO during execution of command: '$BASH_COMMAND'. Exiting with status $?. Output: 'failed'" > ${DATA_DIR}/$LOG_FILE
    exit 1
}

# Set the error trap
trap handle_error ERR

function clear_generated_folder() {
    rm -rf ${GENERATED_DIR}/*
    echo "Cleared the generated folder for new execution!"
}

function step1_fastqc_trimmomatic() {
    fastqc ${DATA_DIR}/${FASTQ_FILE} -o ${GENERATED_DIR}
    trimmomatic SE -threads 4 ${DATA_DIR}/${FASTQ_FILE} ${GENERATED_DIR}/${TRIMMED_FILE} TRAILING:10 -phred33
    echo "Trimmomatic finished running!"
    fastqc ${GENERATED_DIR}/${TRIMMED_FILE} -o ${GENERATED_DIR}
}

function step2_hisat2() {
    mkdir -p ${GENERATED_DIR}/HISAT2
    hisat2 -q --rna-strandness R -x ${WORKFLOW_DIR}/HISAT2/grch38/genome -U ${GENERATED_DIR}/${TRIMMED_FILE} | samtools sort -o ${GENERATED_DIR}/HISAT2/${TRIMMED_BAM}
    echo "HISAT2 finished running!"
}

function step3_featureCounts() {
    mkdir -p ${GENERATED_DIR}/quants
    featureCounts -a ${WORKFLOW_DIR}/hg38/Homo_sapiens.GRCh38.106.gtf -o ${GENERATED_DIR}/quants/${FEATURE_COUNTS} ${GENERATED_DIR}/HISAT2/${TRIMMED_BAM}
    echo "featureCounts finished running!"
}

function print_duration() {
    duration=$SECONDS
    echo "$(($duration / 60)) minutes and $(($duration % 60)) seconds elapsed."
    echo "$duration" > ${GENERATED_DIR}/$DURATION_FILE
    echo "$duration" >  ${DATA_DIR}/$LOG_FILE
}

# Execute functions
clear_generated_folder

# Create files if they do not exist
touch ${GENERATED_DIR}/${DURATION_FILE}

step1_fastqc_trimmomatic
step2_hisat2
step3_featureCounts
print_duration
