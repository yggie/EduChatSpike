// to be moved to a separate file
var Client = {
  connection: null,
  room: null,
  joined: false,
  participants: null,
  nickname: 'my-nickname-other',

  log: function(msg) {
    $('#chat').append("<p>" + msg + "</p>");
  },

  send_ping: function(to) {
    var ping = $iq({
      to: to,
      type: "get",
      id: "ping1"
    }).c("ping", { xmlns: "urn:xmpp:ping" });

    Client.connection.send(ping);
  },

  show_traffic: function(body, type) {
    if (body.childNodes.length > 0) {
      var console = $('#console').get(0);
      var at_bottom = console.scrollTop >= console.scrollHeight - console.clientHeight;

      $.each(body.childNodes, function() {
        $('#console').append('<div class="separator" /><div class="' + type + '">' + Client.pretty_xml(this) + '</div>');
      });

      if (at_bottom) {
        console.scrollTop = console.scrollHeight;
      }
    }
  },

  pretty_xml: function(xml, level) {
    var i, j;
    var result = [];
    if (!level) {
      level = 0;
    }

    result.push("<div class='xml_level" + level + "'>");
    result.push("<span class='xml_punc'>&lt;</span>");
    result.push("<span class='xml_tag'>");
    result.push(xml.tagName);
    result.push("</span>");

    // attributes
    var attrs = xml.attributes;
    var attr_lead = [];
    for (i = 0; i < xml.tagName.length + 1; i++) {
      attr_lead.push("&nbsp;");
    }

    attr_lead = attr_lead.join("");

    for (i = 0; i < attrs.length; i++) {
      result.push(" <span class='xml_aname'>");
      result.push(attrs[i].nodeName);
      result.push("</span><span class='xml_punc'>='</span>");
      result.push("<span class='xml_avalue'>");
      result.push(attrs[i].nodeValue);
      result.push("</span><span class='xml_punc'>'</span>");

      if (i != attrs.length - 1) {
        result.push("</div><div class='xml_level" + level + "'>");
        result.push(attr_lead);
      }
    }

    if (xml.childNodes.length === 0) {
      result.push("<span class='xml_punc'>/&gt;</span></div>");
    } else {
      result.push("<span class='xml_punc'>&gt;</span></div>");

      // children
      $.each(xml.childNodes, function() {
        if (this.nodeType === 1) {
          result.push(Client.pretty_xml(this, level + 1));
        } else if (this.nodeType === 3) {
          result.push("<div class='xml_text xml_level" + (level + 1) + "'>");
          result.push(this.nodeValue);
          result.push("</div>");
        }
      });

      result.push("<div class='xml xml_level" + level + "'>");
      result.push("<span class='xml_punc'>&lt;/</span>");
      result.push("<span class='xml_tag'>");
      result.push(xml.tagName);
      result.push("</span>");
      result.push("<span class='xml_punc'>&gt;</span></div>");
    }

    return result.join("");
  },

  text_to_xml: function(text) {
    var doc = null;
    if (window['DOMParser']) {
      var parser = new DOMParser();
      doc = parser.parseFromString(text, 'text/xml');
    } else if (window['ActiveXObject']) {
      var doc = new ActiveXObject("MSXML2.DOMDocument");
      doc.async = false;
      doc.loadXML(text);
    } else {
      throw {
        type: 'ClientParsingError',
        message: 'No DOMParser object found.'
      };
    }

    var elem = doc.documentElement;
    if ($(elem).filter('parsererror').length > 0) {
      return null;
    }

    return elem;
  },

  on_presence: function(presence) {
    var from = $(presence).attr('from');
    var room = Strophe.getBareJidFromJid(from);

    // make sure this presence is for the right room
    if (room === Client.room) {
      var nick = Strophe.getResourceFromJid(from);

      if ($(presence).attr('type') === 'error' && !Client.joined) {
        // error joining room, reset app
        Client.connection.disconnect();
      } else if (!Client.participants[nick] && $(presence).attr('type') !== 'unavailable') {
        Client.participants[nick] = true;
        $('#participant-list').append('<li>' + nick + '</li>');

        if (Client.joined) {
          $(document).trigger('user_joined', nick);
        }
      } else if (Client.participants[nick] && $(presence).attr('type') === 'unavailable') {
        $('#participant-list li').each(function() {
          if (nick === $(this).text()) {
            $(this).remove();
            return false;
          }
        });

        $(document).trigger('user_left', nick);
      }

      if ($(presence).attr('type') !== 'error' && !Client.joined) {
        // check for status 110 to see if it's our own presence
        if ($(presence).find("status[code='110']").length > 0) {
          // check if server changed our nick
          if ($(presence).find("status[code='210']").length > 0) {
            Client.nickname = Strophe.getResourceFromJid(from);
          }

          // room join complete
          $(document).trigger("room_joined");
        }
      }
    }

    return true;
  },

  on_public_message: function(message) {
    var from = $(message).attr('from');
    var room = Strophe.getBareJidFromJid(from);
    var nick = Strophe.getResourceFromJid(from);

    // make sure message is from the right place
    if (room === Client.room) {
      // is message from a user or the room itself?
      var notice = !nick;

      // messages from ourself will be styled differently
      var nick_class = "nick";
      if (nick === Client.nickname) {
        nick_class += " self";
      }

      var body = $(message).children('body').text();

      var delayed = $(message).children("delay").length > 0 || $(message).children("x[xmlns='jabber:x:delay']").length > 0;

      if (!notice) {
        var delay_css = delayed ? " delayed" : "";
        Client.add_message("<div class='message" + delay_css + "'>&lt;<span class='" + nick_class + "'>" + nick + "</span>&gt; <span class='body'>" + body + "</span></div>");
      } else {
        Client.add_message("<div class='notice'>*** " + body + "</div>");
      }
    }

    return true;
  },

  add_message: function(msg) {
    var chat = $('#chat').get(0);
    var at_bottom = chat.scrollTop >= chat.scrollHeight - chat.clientHeight;

    $('#chat').append(msg);

    if (at_bottom) {
      chat.scrollTop = chat.scrollHeight;
    }
  }
};

(function() {
  function initialize() {
    // var login = { jid: 'testuser@examples.org', password: 'embeddedchatforall', host: 'http://localhost:3000/http-bind', room: 'educhattestroom@educhat.spike' }
    var login = { jid: 'educhatspiketest@blah.im', password: 'randomrandom', host: 'http://bosh.metajack.im:5280/xmpp-httpbind', room: 'educhat-spike-room@rooms.blah.im'}

    Client.room = login.room;
    var conn = new Strophe.Connection(login.host);
    conn.xmlInput = function(body) {
      Client.show_traffic(body, 'incoming');
    };
    conn.xmlOutput = function(body) {
      Client.show_traffic(body, 'outgoing');
    };
    conn.connect(login.jid, login.password, function(status) {
      switch (status) {
        case Strophe.Status.CONNECTED:
          $(document).trigger('connected');
          break;

        case Strophe.Status.DISCONNECTED:
          $(document).trigger('disconnected');
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

      $('.button').removeAttr('disabled');
      $('#input').removeClass('disabled').removeAttr('disabled');

      Client.joined = false;
      Client.participants = {};
      Client.connection.send($pres().c('priority').t('-1'));
      Client.connection.addHandler(Client.on_presence, null, "presence");
      Client.connection.send($pres({ to: Client.room + '/' + Client.nickname}).c('x', {xmlns: 'http://jabber.org/protocol/muc'}));
      Client.connection.addHandler(Client.on_public_message, null, "message", "groupchat");

      var domain = Strophe.getDomainFromJid(Client.connection.jid);

      Client.send_ping(domain);
    });

    $(document).bind('disconnected', function() {
      Client.log("Connection terminated.");

      $('.button').attr('disabled', '');
      $('#input').addClass('disabled').attr('disabled', '');
      $('#participant-list').empty();
      $('#room-name').empty();
      $('#room-topic').empty();
      $('#chat').empty();

      Client.connection = null;
    });

    $(document).bind('authfail', function() {
      Client.log("Authentication failed!");
    });

    $(document).bind("connfail", function() {
      Client.log("Connection failed!");
    });

    $(document).bind('room_joined', function() {
      Client.joined = true;

      $('#leave').removeAttr('disabled');
      $('#room-name').text(Client.room);

      $('#chat').append("<div class='notice'>*** Room joined.</div>");
    });

    $(document).bind('user_joined', function(ev, nick) {
      Client.add_message("<div class='notice'> *** " + nick + " joined.</div>");
    });

    $(document).bind('user_left', function(ev, nick) {
      Client.add_message("<div class='notice'>*** " + nick + " left.</div>");
    });

    $('#leave').click(function() {
      Client.connection.send($pres({ to: Client.room + '/' + Client.nickname, type: 'unavailable' }));
      Client.connection.disconnect();
    });

    $('#send_button').click(function() {
      var input = $('#input').val();
      var error = false;
      if (input.length > 0) {
        if (input[0] === '<') {
          var xml = Client.text_to_xml(input);
          if (xml) {
            Client.connection.send(xml);
            $('#input').val('');
          } else {
            error = true;
          }
        } else if (input[0] === '$') {
          try {
            var builder = eval(input);
            Peek.connection.send(builder);
            $('#input').val('');
          } catch (e) {
            error = true
          }
        } else {
          error = true;
        }
      }

      if (error) {
        $('#input').animate({ backgroundColor: "#faa" });
      }
    });

    $('#input').keypress(function() {
      $(this).css({ backgroundColor: "#fff" });
    });

    $('#chat_input').keypress(function(ev) {
      if (ev.which === 13) {
        ev.preventDefault();

        var body = $(this).val();

        Client.connection.send($msg({ to: Client.room, type: "groupchat" }).c('body').t(body));

        $(this).val('');
      }
    });
  }

  $(document).ready(initialize);
})();

