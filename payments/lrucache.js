const LRU = require("lru-cache")


const options = {
    max: 50000
    , maxAge: 1000 * 60 * 60
}

let cache = new LRU(options)



exports.set = (key, value) => {
    return cache.set(key, value)
}

exports.get = (key) => {
    console.log(cache.dump())
    return cache.get(key)
}