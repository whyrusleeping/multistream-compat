var net = require('net')
var Select = require('multistream-select').Select
var argv = require('minimist')(process.argv.slice(2))

addr = argv.addr || 5050
proto = argv.protos || "/test"

proto = proto.split(',')


var ms = new Select()
function echo(sock) {
	sock.pipe(sock)
}

for (var p in proto) {
	ms.addHandler(proto[p], echo)
}

var server = net.createServer(function(s) {
	ms.handle(s)
}).listen(addr)
