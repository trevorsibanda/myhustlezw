var express = require('express');
var router = express.Router();

router.get('/', function(req, res, next) {
  res.send('respond with a resource');
});


//ecocash callback
router.get('/ecocash', function (req, res, next) {
  console.log('request details:')
  console.log('cookies:', req.cookies)
  console.log('headers:', req.headers)
  console.log('query:', req.query)
  console.log('body:', req.body.toString())
  console.log('params:', req.params)

  res.json({status: 'ok'})
})

//ecocash callback
router.post('/ecocash', function (req, res, next) {
  console.log('request details:')
  console.log('cookies:', req.cookies)
  console.log('headers:', req.headers)
  console.log('query:', req.query)
  console.log('body:', req.body.toString())
  console.log('params:', req.params)

  res.json({status: 'ok'})
});

module.exports = router;
