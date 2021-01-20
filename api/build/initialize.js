"use strict";

var _Object$defineProperty = require("@babel/runtime-corejs3/core-js/object/define-property");

_Object$defineProperty(exports, "__esModule", {
  value: true
});

exports.initializeDb = void 0;

const initializeDb = driver => {
  const initCyper = `call apoc.schema.assert({}, {Movie:["movieId"], User:["userId"]})`;

  const executeQuery = driver => {
    const session = driver.session();
    return session.writeTransaction(tx => tx.run(initCyper)).finally(() => {
      session.close();
      console.log("connected to db succesfully!");
    }).catch(error => {
      console.log("connection to db could not be established!\n", error.message);
    });
  };

  executeQuery(driver);
};

exports.initializeDb = initializeDb;