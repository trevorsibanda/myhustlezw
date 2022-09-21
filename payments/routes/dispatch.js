var express = require('express');
const ecocash = require('../ecocash')
const LRU = require("../lrucache")
const m = require('./middleware')
var router = express.Router();

/* GET users listing. */
router.get('/', function(req, res, next) {
  res.send('respond with a resource');
});


//ecocash dispatch
router.post('/ecocash', function (req, res, next) {
  if (!m.isAuthorized(req, res)) {
    return
  }
  ecocash.init(req.body.national_number, req.body.price, req.body._id, 'MyHustleZW ', req.body.email).then(function (response) {
    if (response.success) {
      // These are the instructions to show the user. 
      // Instruction for how the user can make payment
      LRU.set(req.body._id, response)
    } else {
      console.log('Encountered error: ', response.error)
    }

    res.send(response)
  }).catch(ex => {
    // Ahhhhhhhhhhhhhhh
    // *freak out*
    console.log('Failed with exception: ', ex)
    res.send({error: 'Internal server error', details: ex})
  })
})

//2checkout dispatch
router.post('/', function (req, res, next) {
  res.send('respond with a resource');
});

module.exports = router;
