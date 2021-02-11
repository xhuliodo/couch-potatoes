FROM neo4j:4.2.0

RUN apt-get update; apt-get install -y curl

ENV NEO4J_AUTH=neo4j/letmein \
    APOC_VERSION=4.2.0.0 

ENV APOC_URI https://github.com/neo4j-contrib/neo4j-apoc-procedures/releases/download/${APOC_VERSION}/apoc-${APOC_VERSION}-all.jar
RUN sh -c 'cd /var/lib/neo4j/plugins && curl -L -O "${APOC_URI}"'

EXPOSE 7474 7473 7687

ENV IMPORT_ENTRYPOINT_URI  https://raw.githubusercontent.com/xhuliodo/couch-potatoes/master/db/import/import-db-entrypoint.sh

RUN sh -c 'cd /var/lib/neo4j/ && curl -L -O "${IMPORT_ENTRYPOINT_URI}"'

RUN chmod -R 777 /var/lib/neo4j/import-db-entrypoint.sh

ENTRYPOINT [ "./import-db-entrypoint.sh" ]