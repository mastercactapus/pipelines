FROM mastercactapus/build-env:stretch-amd64
RUN apt-get -qq update && \
    apt-get -qq install -y libc6-dev libssl-dev libexpat1-dev libavcodec-dev libgl1-mesa-dev libqt4-dev && \
    rm -rf /var/lib/apt/lists/*
