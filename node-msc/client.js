var tcp = require('net')
var abr = require('async-buffered-reader')
var crypto = require('crypto');
var Interactive = require('multistream-select').Interactive
var argv = require('minimist')(process.argv.slice(2))

addr = argv.addr || 5050
proto = argv.protos || "/test"

var socket = tcp.connect({port: addr}, connected)

function random (howMany, chars) {
    chars = chars 
        || "abcdefghijklmnopqrstuwxyzABCDEFGHIJKLMNOPQRSTUWXYZ0123456789";
    var rnd = crypto.randomBytes(howMany)
        , value = new Array(howMany)
        , len = chars.length;

    for (var i = 0; i < howMany; i++) {
        value[i] = chars[rnd[i] % len]
    };

    return value.join('');
}

function connected () {
  var msi = new Interactive()

  msi.handle(socket, function () {
      msi.select(proto, function (err, ds) {
        if (err) {
          return console.log(err)
		}

		stuff = random(4096)

		ds.write(stuff)
		abr(ds, 4096, function(out) {
			if (out.toString() !== stuff.toString()) {
				console.log(out.toString())
				console.log("data was incorrect")
				ds.end()
				return
			}
			console.log("data was correct")
			ds.end()
		})
	  })
  })
}
