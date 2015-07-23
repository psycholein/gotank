// Generated by CoffeeScript 1.9.3
(function() {
  var Compass,
    bind = function(fn, me){ return function(){ return fn.apply(me, arguments); }; },
    extend = function(child, parent) { for (var key in parent) { if (hasProp.call(parent, key)) child[key] = parent[key]; } function ctor() { this.constructor = child; } ctor.prototype = parent.prototype; child.prototype = new ctor(); child.__super__ = parent.prototype; return child; },
    hasProp = {}.hasOwnProperty;

  window.Compass = Compass = (function(superClass) {
    extend(Compass, superClass);

    function Compass() {
      this.router = bind(this.router, this);
      this.config = bind(this.config, this);
      return Compass.__super__.constructor.apply(this, arguments);
    }

    Compass.prototype.config = function() {
      return this.position = "left";
    };

    Compass.prototype.router = function(event) {
      var $module;
      switch (event.Task) {
        case "compass":
          $module = $(this.selector).filter("[data-name='" + event.Name + "']");
          return $module.find('.value').text(event.Data.value);
      }
    };

    return Compass;

  })(window.BasicModule);

}).call(this);
