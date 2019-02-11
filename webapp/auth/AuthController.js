var express = require('express');
var router = express.Router();
var bodyParser = require('body-parser');
router.use(bodyParser.urlencoded({ extended: false }));
router.use(bodyParser.json());
var User = require('../User');

var jwt = require('jsonwebtoken');
var bcrypt = require('bcryptjs');
var config = require('../config');

// https://www.toptal.com/nodejs/secure-rest-api-in-nodejs
router.post('/register', function(req, res) {
  
    var hashedPassword = bcrypt.hashSync(req.body.password, 8);
    
    User.create({
      name : req.body.name,
      email : req.body.email,
      password : hashedPassword
    },
    function (err, user) {
      if (err) return res.status(500).send("There was a problem registering the user.")
      // create a token
      var token = jwt.sign({ id: user._id }, config.secret, {
        expiresIn: 86400 // expires in 24 hours
      });
      res.status(200).send({ auth: true, token: token });
    }); 
  });

router.post('/login', function(req, res) {


});

router.get('/authenticate', function(req, res) {
    var token = req.headers['x-access-token'];
    if (!token) return res.status(401).send({ auth: false, message: 'No token provided.' });
    jwt.verify(token, config.secret, function(err, decoded) {
        if (err) return res.status(500).send({ auth: false, message: 'Failed to authenticate token.' }); 
        User.findById(decoded.id, function (err, user) {
            if (err) return res.status(500).send("There was a problem finding the user.");
            if (!user) return res.status(404).send("No user found.");
            res.status(200).send(decoded);
        });
    });
});

exports.isPasswordAndUserMatch = (req, res, next) => {
    User.findByUsername(req.body.username)
        .then((user)=>{
            if(!user[0]){
                res.status(404).send({});
            }else{
                const passwordHash = user[0].password;
                
                if (hash === passwordFields[1]) {
                    return next();
                } else {
                    return res.status(401).send({errors: ['Invalid email or password']});
                }
            }
        });
 };


module.exports = router;