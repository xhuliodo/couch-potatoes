export const initializeDb = (driver) => {
  const initCyper = `call apoc.schema.assert({}, {Movie:["movieId"], User:["userId"]}, false)`;

  const executeQuery = (driver) => {
    const session = driver.session();
    return session
      .writeTransaction((tx) => tx.run(initCyper))
      .finally(() => {
        session.close();
        console.log("connected to db succesfully!");
      })
      .catch((error) => {
        console.log(
          "connection to db could not be established!\n",
          error.message
        );
      });
  };

  executeQuery(driver);
};
