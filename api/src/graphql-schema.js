import fs from "fs";
import path from "path";

// imports written schema
export const typeDefs = fs
  .readFileSync(path.join(__dirname, "schema.graphql"))
  .toString("utf-8");
