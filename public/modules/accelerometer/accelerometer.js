// Generated by CoffeeScript 1.9.3
(function() {
  var Accelerometer,
    bind = function(fn, me){ return function(){ return fn.apply(me, arguments); }; },
    extend = function(child, parent) { for (var key in parent) { if (hasProp.call(parent, key)) child[key] = parent[key]; } function ctor() { this.constructor = child; } ctor.prototype = parent.prototype; child.prototype = new ctor(); child.__super__ = parent.prototype; return child; },
    hasProp = {}.hasOwnProperty;

  window.Accelerometer = Accelerometer = (function(superClass) {
    extend(Accelerometer, superClass);

    function Accelerometer() {
      this.router = bind(this.router, this);
      this.config = bind(this.config, this);
      return Accelerometer.__super__.constructor.apply(this, arguments);
    }

    Accelerometer.prototype.config = function() {
      return this.position = "right";
    };

    Accelerometer.prototype.router = function(event) {
      var $module, results, value;
      switch (event.Task) {
        case "accelerometer":
          $module = $(this.selector).filter("[data-name='" + event.Name + "']");
          results = [];
          for (value in event.Data) {
            results.push($module.find("." + value).text(event.Data[value]));
          }
          return results;
      }
    };

    return Accelerometer;

  })(window.BasicModule);

}).call(this);
