var createError = require('http-errors');
var express = require('express');
var path = require('path');
var logger = require('morgan');
var bodyParser = require('body-parser');
var serverless = require('serverless-http');

var pollStatusRouter = require('./routes/poll');
var dispatchRouter = require('./routes/dispatch');
var callbackRouter = require('./routes/callback');

var app = express();

// view engine setup
app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'jade');

app.use(logger('dev'));
app.use(bodyParser.json());
//app.use(express.urlencoded({ extended: false }));

app.use('/poll', pollStatusRouter);
app.use('/dispatch', dispatchRouter);
app.use('/callback', callbackRouter);

// catch 404 and forward to error handler
app.use(function(req, res, next) {
  next(createError(404));
});

// error handler
app.use(function(err, req, res, next) {
  // set locals, only providing error in development
  res.locals.message = err.message;
  res.locals.error = req.app.get('env') === 'development' ? err : {};

  // render the error page
  res.status(err.status || 500);
  res.render('error');
});

module.exports.handler = serverless(app);
