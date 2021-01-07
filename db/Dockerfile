FROM neo4j:4.2.1

RUN apt-get update; apt-get install -y curl

ENV NEO4J_AUTH=neo4j/letmein \
    NEO4J_PASSWD=letmein \
    APOC_VERSION=4.2.0.0 

COPY ./import /var/lib/neo4j/import/

CMD bin/neo4j-admin set-initial-password ${NEO4J_PASSWD} || true && \
    bin/neo4j start && sleep 5 && \
    cat /var/lib/neo4j/import/import.cypher | NEO4J_USERNAME=neo4j NEO4J_PASSWORD=${NEO4J_PASSWD} bin/cypher-shell --fail-fast; \
    done && \
    tail -f /var/lib/neo4j/logs/neo4j.log

# ENV APOC_URI https://github.com/neo4j-contrib/neo4j-apoc-procedures/releases/download/${APOC_VERSION}/apoc-${APOC_VERSION}-all.jar
# RUN sh -c 'cd /var/lib/neo4j/plugins && curl -L -O "${APOC_URI}"'

EXPOSE 7474 7473 7687

CMD ["neo4j"]