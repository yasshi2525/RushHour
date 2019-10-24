FROM python as base

WORKDIR /data

COPY src src

RUN apt-get -y update && apt-get install fonts-vlgothic && \
    pip install Sphinx \ 
        sphinxcontrib-blockdiag \
        sphinxcontrib-seqdiag \
        sphinxcontrib-actdiag && \
    sphinx-build -b html src dist

FROM nginx
COPY --from=base /data/dist /usr/share/nginx/html