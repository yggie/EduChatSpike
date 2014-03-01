var xmpp = require('node-xmpp-server')
    http = require('http');

var sv = new xmpp.BOSHServer({
  server: http.createServer(function(req, res) {
    console.log(req);
    sv.handleHTTP(req, res);
  }).listen(9000)
});
