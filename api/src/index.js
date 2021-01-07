import express from "express";
import { ApolloServer } from "apollo-server-express";
import neo4j from "neo4j-driver";
import { makeAugmentedSchema } from "neo4j-graphql-js";
import dotenv from "dotenv";

import { typeDefs } from "./graphql-schema";
import { initializeDb } from "./initialize";

// set env vars from .env file
dotenv.config();

const app = express();

// creates autogenerated queries and mutations from ./schema.graphql
// for types that you dont want to do that, you can exclude them
const schema = makeAugmentedSchema({
  typeDefs,
  config: {
    experimental: true,
  },
});

// creates a connection instance with credentials provided in .env file
const driver = neo4j.driver(
  process.env.NEO4J_URI || "bolt://localhost:7687",
  neo4j.auth.basic(
    process.env.NEO4J_USER || "neo4j",
    process.env.NEO4J_PASSWORD || "letmein"
  ),
  {
    encrypted: process.env.NEO4J_ENCRYPTED ? "ENCRYPTION_ON" : "ENCRYPTION_OFF",
  }
);

// performing any initialization procedure necessary
// (ex: creating constraints, ensuring indexes are online, etc...)
const init = async (driver) => {
  await initializeDb(driver);
};

init(driver);

// creates an apollo server instances with the generated schema from
// makeAugmentedSchema. it also injects the neo4j driver instance
// into the context making it possible for the auto generated
// resolvers to connect to the database
// TODO: set introspection and playground to false in production
const server = new ApolloServer({
  context: { driver, neo4jDatabase: process.env.NEO4J_DATABASE },
  schema,
  introspection: true,
  playground: true,
});

// define endpoint details
const port = process.env.GRAPHQL_SERVER_PORT || 4001;
const path = process.env.GRAPHQL_SERVER_PATH || "/graphql";
const host = process.env.GRAPHQL_SERVER_HOST || "0.0.0.0";

// adding apollo instance as a middleware to express instance
// TODO: apply auth
server.applyMiddleware({ app, path });

// activate express instance
app.listen({ host, port, path }, () => {
  console.log(
    `couch-potatoes api is listening on http://${host}:${port}${path}`
  );
});
