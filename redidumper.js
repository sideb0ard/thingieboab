#!/usr/bin/env node
var redis = require("redis");
var language = require("./language.json");

var client = redis.createClient();
client.on("error", function (err) {
  console.log("Error " + err);
});
var keysz = Object.keys(language);
for (var i = 0; i < keysz.length; i++) {
  var componentz = Object.getOwnPropertyNames(language[keysz[i]]);
  for (var j = 0; j < componentz.length; j++) {
    var rk = "aigor:" + keysz[i] + ":" + componentz[j];
    var rv = language[keysz[i]][componentz[j]];
    console.log("Setting redis value for:" + rk + " TO: " + rv);
    client.set(rk, rv, redis.print);
  }
}


client.quit();
