// Generated by CoffeeScript 1.9.3
(function() {
  var App, Event, Network, Resources, root, test,
    bind = function(fn, me){ return function(){ return fn.apply(me, arguments); }; },
    indexOf = [].indexOf || function(item) { for (var i = 0, l = this.length; i < l; i++) { if (i in this && this[i] === item) return i; } return -1; };

  root = typeof exports !== "undefined" && exports !== null ? exports : this;

  root.App = App = (function() {
    function App() {
      var url;
      url = "ws://" + document.location.host + "/ws";
      this.network = new Network(url);
      this.event = new Event(this.network);
      this.resources = new Resources(this.event);
      this.network.connect(this.event);
    }

    return App;

  })();

  root.Network = Network = (function() {
    function Network(url1) {
      this.url = url1;
      this.send = bind(this.send, this);
      this.disconnect = bind(this.disconnect, this);
      this.connect = bind(this.connect, this);
      this.ws = null;
    }

    Network.prototype.connect = function(event) {
      this.ws = new WebSocket(this.url);
      this.ws.onopen = event.connected;
      this.ws.onclose = event.disconnected;
      this.ws.onmessage = event.receive;
      this.ws.onerror = event.error;
      return this.event = event;
    };

    Network.prototype.disconnect = function() {
      this.ws.close();
      return this.ws = null;
    };

    Network.prototype.send = function(msg) {
      if (this.ws && msg) {
        return this.ws.send(msg);
      }
    };

    return Network;

  })();

  Event = (function() {
    function Event(network) {
      this.network = network;
      this.send = bind(this.send, this);
      this.receive = bind(this.receive, this);
      this.register = bind(this.register, this);
      this.registered = {};
    }

    Event.prototype.register = function(srcModule, destModuleFunc) {
      if (!this.registered[srcModule]) {
        this.registered[srcModule] = [];
      }
      return this.registered[srcModule].push(destModuleFunc);
    };

    Event.prototype.connected = function(e) {};

    Event.prototype.disconnected = function(e) {};

    Event.prototype.receive = function(e) {
      var data, func, i, j, len, len1, ref, ref1, results;
      data = JSON.parse(e.data);
      if (!data) {
        return;
      }
      if (this.registered[data.Module]) {
        ref = this.registered[data.Module];
        for (i = 0, len = ref.length; i < len; i++) {
          func = ref[i];
          func(data);
        }
      }
      if (this.registered['_all']) {
        ref1 = this.registered['_all'];
        results = [];
        for (j = 0, len1 = ref1.length; j < len1; j++) {
          func = ref1[j];
          results.push(func(data));
        }
        return results;
      }
    };

    Event.prototype.error = function(e) {};

    Event.prototype.send = function(module, name, task, value) {
      var data;
      data = {
        Module: module,
        Name: name,
        Task: task,
        Value: value
      };
      return this.network.send(JSON.stringify(data));
    };

    return Event;

  })();

  root.Resources = Resources = (function() {
    function Resources(event1) {
      this.event = event1;
      this.loadResource = bind(this.loadResource, this);
      this.initModule = bind(this.initModule, this);
      this.loadModule = bind(this.loadModule, this);
      this.router = bind(this.router, this);
      this.modules = {};
      this.event.register('_all', this.router);
    }

    Resources.prototype.router = function(event) {
      if (!(event.Name === 'module' || event.Task === 'web')) {
        return;
      }
      switch (event.Value) {
        case "load":
          return this.loadModule(event);
        case "init":
          return this.initModule(event);
      }
    };

    Resources.prototype.loadModule = function(event) {
      var css, file, ref;
      if (ref = event.Name, indexOf.call(this.modules, ref) >= 0) {
        return;
      }
      this.modules[event.Module] = {};
      file = "modules/" + event.Module + "/" + event.Module;
      this.loadResource(event.Module, file + ".js", 'script');
      this.loadResource(event.Module, file + ".ect", 'text');
      css = $("<link rel='stylesheet' href='" + file + ".css' type='text/css' />");
      return $("head").append(css);
    };

    Resources.prototype.initModule = function(event) {
      console.log("init", event);
      if (event.Module) {
        return new event.Module(event.Name);
      }
    };

    Resources.prototype.loadResource = function(module, file, type) {
      return $.ajax({
        url: file,
        dataType: type,
        success: (function(_this) {
          return function(data) {
            return _this.modules[module][type] = data;
          };
        })(this)
      });
    };

    return Resources;

  })();

  test = "App";

  $(function() {
    return new root[test];
  });

}).call(this);
