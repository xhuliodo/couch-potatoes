FROM neo4j:4.2.1-enterprise

RUN apt-get update; apt-get install -y curl

ENV APOC_VERSION=4.2.0.0 \
    NEO4J_ACCEPT_LICENSE_AGREEMENT=yes

ENV APOC_URI https://github.com/neo4j-contrib/neo4j-apoc-procedures/releases/download/${APOC_VERSION}/apoc-${APOC_VERSION}-all.jar
RUN sh -c 'cd /var/lib/neo4j/plugins && curl -L -O "${APOC_URI}"'

EXPOSE 7474 7473 7687

CMD ["neo4j"]