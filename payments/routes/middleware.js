const sha512 = require('js-sha512');

function isAuthorized(req, res, next) {
    return true;
}

module.exports.isAuthorized = isAuthorized;