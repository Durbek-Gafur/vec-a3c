FROM golang:latest 

ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOPROXY=https://proxy.golang.org,direct \
    GONOPROXY=none \
    GOSUMDB=sum.golang.org \
    GONOSUMDB=none \ 
    PATH="/miniconda/bin:${PATH}"

WORKDIR /app
COPY ./ /app

RUN mkdir "/build"

RUN go mod tidy

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        build-essential \
        default-jre \
        fastqc \
        samtools \
        wget \
        gzip \
        tar \
        unzip \
        zlib1g-dev && \
    rm -rf /var/lib/apt/lists/* && \
    wget https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh && \
    bash Miniconda3-latest-Linux-x86_64.sh -b -p /miniconda && \
    rm Miniconda3-latest-Linux-x86_64.sh 

COPY workflow/environment.yml .

# Create and activate the conda environment, make the conda environment available when the container is run
RUN conda env create -f environment.yml && \
    echo "source activate $(head -1 environment.yml | cut -d' ' -f2)" > ~/.bashrc && \
    PATH="/miniconda/envs/$(head -1 environment.yml | cut -d' ' -f2)/bin:$PATH" && \
    conda clean -a -y


# Make the bioinformatics rna script executable 
RUN chmod +x ${PROJECT_DIR}/workflow/rna.sh 

RUN go build -o /build/app ./cmd/vec_worker

ENTRYPOINT ["/build/app"]
