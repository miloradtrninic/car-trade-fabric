var express = require("express");
var app = express();
var network = require('./network');
var db = require('./db');

var UserController = require('./UserController');
var AuthController = require('./auth/AuthController');

app.use('/users', UserController);
app.use('/auth', AuthController);

app.get("/test", (req, res, next) => {
    network.testFunction(req, res, next);
});

app.put("/changeColor", (req, res, next) => AuthController.isTokenValid(req, res, next), (req, res, next) => {
    network.changeColor(req, res, next);
});

app.get("/getByColor", (req, res, next) => AuthController.isTokenValid(req, res, next), (req, res, next) => {
    network.getByColor(req, res, next);
});

app.get("/getByColor", (req, res, next) => AuthController.isTokenValid(req, res, next), (req, res, next) => {
    network.getByColor(req, res, next);
});

app.get("/getCarHistory", (req, res, next) => AuthController.isTokenValid(req, res, next), (req, res, next) => {
    network.getCarHistory(req, res, next);
});

app.put("/changeOwner", (req, res, next) => AuthController.isTokenValid(req, res, next), (req, res, next) => {
    network.changeOwner(req, res, next);
});

app.put("/noteDamage", (req, res, next) => AuthController.isTokenValid(req, res, next), (req, res, next) => {
    network.noteDamage(req, res, next);
});

app.put("/repairDamage", (req, res, next) => AuthController.isTokenValid(req, res, next), (req, res, next) => {
    network.noteDamage(req, res, next);
});

app.put("/initLedger", (req, res, next) => AuthController.isTokenValid(req, res, next), (req, res, next) => {
    network.initLedger(req, res, next);
});

app.listen(3000, () => {
 console.log("Server running on port 3000");
});



