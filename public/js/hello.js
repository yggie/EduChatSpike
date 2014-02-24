(function() {
  var Client = {
    connection: null,
    log: function(msg) {
      $('#log').append("<p>" + msg + "</p>");
    }
  };

  var conn = new Strophe.Connection("/http-bind");
  conn.connect("thisisanexample@examples.org", null, function(status) {
    switch (status) {
      case Strophe.Status.CONNECTED:
        $(document).trigger('connected');
        break;

      case Strophe.Status.CONNECTING:
        $(document).trigger("connecting");
        break;

      case Strophe.Status.DISCONNECTED:
        $(document).trigger('disconnected');
        break;

      case Strophe.Status.DISCONNECTING:
        $(document).trigger("disconnecting");
        break;

      case Strophe.Status.CONNFAIL:
        $(document).trigger('connfail');
        break;

      case Strophe.Status.AUTHFAIL:
        $(document).trigger('authfail');
        break;
    }

    Client.connection = conn;
  });

  $(document).bind('connected', function() {
    Client.log("Connection established.");
  });

  $(document).bind("connecting", function() {
    Client.log("Connecting...");
  });

  $(document).bind('disconnected', function() {
    Client.log("Connection terminated.");

    Client.connection = null;
  });

  $(document).bind('disconnecting', function() {
    Client.log("Disconnecting from server.");
  });

  $(document).bind('authfail', function() {
    Client.log("Authentication failed!");
  });

  $(document).bind("connfail", function() {
    Client.log("Connection failed!");
  });
})();
