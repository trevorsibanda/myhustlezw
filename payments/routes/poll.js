var express = require('express');
const Ecocash = require("../ecocash")
const LRU = require("../lrucache")
const m = require('./middleware')
const query = require("querystring")

var router = express.Router();

/* poll ecocash transaction */
router.post('/ecocash', function (req, res, next) {
  if (!m.isAuthorized(req, res)) {
    return
  }
  console.log(req.body)
  Ecocash.poll(req.body._id).then(data => {
    let m = query.decode(data)
    res.send(m)
  }).catch(err => {
    res.send({error: err})
  })
});

module.exports = router;
