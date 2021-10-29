#!/bin/bash

docker build -t glab_dev -f dev/Dockerfile .
docker run -it -v $(pwd):/lab glab_dev /bin/bash