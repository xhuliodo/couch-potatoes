"use strict";

var _interopRequireDefault = require("@babel/runtime-corejs3/helpers/interopRequireDefault");

var _Object$defineProperty = require("@babel/runtime-corejs3/core-js/object/define-property");

_Object$defineProperty(exports, "__esModule", {
  value: true
});

exports.typeDefs = void 0;

var _fs = _interopRequireDefault(require("fs"));

var _path = _interopRequireDefault(require("path"));

// imports written schema
const typeDefs = _fs.default.readFileSync(_path.default.join(__dirname, "schema.graphql")).toString("utf-8");

exports.typeDefs = typeDefs;