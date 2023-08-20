#!/bin/bash
source /miniconda/etc/profile.d/conda.sh && conda activate rnaseq

set -e
set -o pipefail

SECONDS=0

WORKFLOW_DIR="./workflow"
DATA_DIR="${WORKFLOW_DIR}/data"
# /app/workflow/data/generated
GENERATED_DIR="${DATA_DIR}/generated"

# Check if a parameter was provided
if [[ -z "$1" ]]; then
    echo "$(date) - Error: Please provide the name of the FASTQ file as a parameter."
    echo "$(date) - Error: Please provide the name of the FASTQ file as a parameter." >> ${DATA_DIR}/$LOG_FILE
    exit 1
fi

FASTQ_FILE="$1"
if [[ ! -f "${DATA_DIR}/${FASTQ_FILE}" ]]; then
    echo "$(date) - Error: File ${DATA_DIR}/${FASTQ_FILE} does not exist." >> ${DATA_DIR}/$LOG_FILE
    exit 1
fi

TRIMMED_FILE="${FASTQ_FILE%.fastq}_trimmed.fastq"
TRIMMED_BAM="${FASTQ_FILE%.fastq}_trimmed.bam"
FEATURE_COUNTS="${FASTQ_FILE%.fastq}_featurecounts.txt"
DURATION_FILE="result.txt"
LOG_FILE="logs.txt"

# Define error trap
function handle_error() {
    echo "Error occurred in the script at line $LINENO during execution of command: '$BASH_COMMAND'. Exiting with status $?. Output: 'failed'" > ${GENERATED_DIR}/$DURATION_FILE
    echo "$(date) - Error occurred in the script at line $LINENO during execution of command: '$BASH_COMMAND'. Exiting with status $?. Output: 'failed'" >> ${DATA_DIR}/$LOG_FILE
    exit 1
}

# Set the error trap
trap handle_error ERR

function clear_generated_folder() {
    rm -rf ${GENERATED_DIR}/*
    echo "$(date) - Cleared the generated folder for new execution!" >> ${DATA_DIR}/$LOG_FILE
}

function step1_fastqc_trimmomatic() {
    fastqc ${DATA_DIR}/${FASTQ_FILE} -o ${GENERATED_DIR}
    echo "$(date) - Finished running fastqc for ${FASTQ_FILE}" >> ${DATA_DIR}/$LOG_FILE
    trimmomatic SE -threads 4 ${DATA_DIR}/${FASTQ_FILE} ${GENERATED_DIR}/${TRIMMED_FILE} TRAILING:10 -phred33
    echo "$(date) - Trimmomatic finished running!" >> ${DATA_DIR}/$LOG_FILE
    fastqc ${GENERATED_DIR}/${TRIMMED_FILE} -o ${GENERATED_DIR}
    echo "$(date) - Finished running fastqc for trimmed file ${TRIMMED_FILE}" >> ${DATA_DIR}/$LOG_FILE
}

function step2_hisat2() {
    mkdir -p ${GENERATED_DIR}/HISAT2

    # Run hisat2 and capture both stdout and stderr
    hisat2 -q --rna-strandness R -x ${WORKFLOW_DIR}/HISAT2/grch38/genome -U ${GENERATED_DIR}/${TRIMMED_FILE} 2> >(tee -a ${DATA_DIR}/$LOG_FILE >> ${GENERATED_DIR}/$DURATION_FILE) | samtools sort -o ${GENERATED_DIR}/HISAT2/${TRIMMED_BAM}

    echo "$(date) - HISAT2 finished running!" >> ${DATA_DIR}/$LOG_FILE
}


function step3_featureCounts() {
    mkdir -p ${GENERATED_DIR}/quants
    featureCounts -a ${WORKFLOW_DIR}/hg38/Homo_sapiens.GRCh38.106.gtf -o ${GENERATED_DIR}/quants/${FEATURE_COUNTS} ${GENERATED_DIR}/HISAT2/${TRIMMED_BAM}
    echo "$(date) - featureCounts finished running!" >> ${DATA_DIR}/$LOG_FILE
}

function print_duration() {
    duration=$SECONDS
    echo "$(date) - $(($duration / 60)) minutes and $(($duration % 60)) seconds elapsed." >> ${DATA_DIR}/$LOG_FILE
    echo "$duration" > ${GENERATED_DIR}/$DURATION_FILE
}

# Execute functions
clear_generated_folder
touch ${GENERATED_DIR}/${DURATION_FILE}

step1_fastqc_trimmomatic
step2_hisat2
step3_featureCounts
print_duration
