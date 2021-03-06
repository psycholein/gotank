// Generated by CoffeeScript 1.9.3
(function() {
  var Motion,
    bind = function(fn, me){ return function(){ return fn.apply(me, arguments); }; },
    extend = function(child, parent) { for (var key in parent) { if (hasProp.call(parent, key)) child[key] = parent[key]; } function ctor() { this.constructor = child; } ctor.prototype = parent.prototype; child.prototype = new ctor(); child.__super__ = parent.prototype; return child; },
    hasProp = {}.hasOwnProperty;

  window.Motion = Motion = (function(superClass) {
    extend(Motion, superClass);

    function Motion() {
      this.router = bind(this.router, this);
      this.config = bind(this.config, this);
      return Motion.__super__.constructor.apply(this, arguments);
    }

    Motion.prototype.config = function() {
      return this.position = "right";
    };

    Motion.prototype.router = function(event) {
      var $module;
      switch (event.Task) {
        case "motion":
          $module = $(this.selector).filter("[data-name='" + event.Name + "']");
          return $module.find('.value').text(event.Data.value);
      }
    };

    return Motion;

  })(window.BasicModule);

}).call(this);
