FROM neo4j:4.4.48

# Credentials are supplied by the environment (NEO4J_AUTH, or NEO4J_USER and
# NEO4J_PASSWORD); the entrypoint falls back to the development defaults.
ENV APOC_VERSION=4.4.0.40

ENV APOC_URI=https://github.com/neo4j-contrib/neo4j-apoc-procedures/releases/download/${APOC_VERSION}/apoc-${APOC_VERSION}-all.jar

RUN wget -q -O /var/lib/neo4j/plugins/apoc-${APOC_VERSION}-all.jar "${APOC_URI}" && \
    chmod 644 /var/lib/neo4j/plugins/apoc-${APOC_VERSION}-all.jar

# The seed files are mounted at /import by docker-compose, so keep that as the
# import directory instead of Neo4j's default (/var/lib/neo4j/import).
ENV NEO4J_dbms_directories_import=/import \
    NEO4J_dbms_security_procedures_unrestricted=apoc.*

# Copy the initialization script from this repository; never fetch it at build time.
COPY import/import-db-entrypoint.sh /import-db-entrypoint.sh
RUN chmod 755 /import-db-entrypoint.sh

EXPOSE 7474 7473 7687

ENTRYPOINT ["/import-db-entrypoint.sh"]
